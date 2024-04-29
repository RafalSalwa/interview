package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/RafalSalwa/auth-api/cmd/tester_service/config"
	"github.com/RafalSalwa/auth-api/pkg/generator"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/models"
)

type Ordered struct {
	ctx                context.Context
	cfg                *config.Config
	client             *http.Client
	logger             *logger.Logger
	endpoint           string
	endpointSignUp     string
	endpointSignIn     string
	endpointAuthCode   string
	endpointVerifyCode string
}

const (
	password                  = "VeryG00dPass!"
	numUsers                  = 20
	maxNbConcurrentGoroutines = 10
)

var (
	concurrentGoroutines = make(chan struct{}, maxNbConcurrentGoroutines)
)

func NewOrdered(ctx context.Context, cfg *config.Config, l *logger.Logger) WorkerRunner {
	ordered := &Ordered{
		ctx:    ctx,
		cfg:    cfg,
		client: &http.Client{},
		logger: l,
	}
	ordered.endpoint = fmt.Sprintf("http://%s", cfg.HTTP.Addr)
	ordered.endpointSignUp = fmt.Sprintf("%s/auth/signup", ordered.endpoint)
	ordered.endpointSignIn = fmt.Sprintf("%s/auth/signin", ordered.endpoint)
	ordered.endpointAuthCode = fmt.Sprintf("%s/auth/code", ordered.endpoint)
	ordered.endpointVerifyCode = fmt.Sprintf("%s/auth/verify", ordered.endpoint)

	return ordered
}

func (o *Ordered) Run() {
	qCreatedUsers := make(chan testUser, numUsers)
	qActivatedUsers := make(chan testUser, numUsers)
	qFailedUsers := make(chan testUser, numUsers)

	done := make(chan bool)

	for i := 1; i <= numUsers; i++ {
		go o.createUser(o.ctx, o.cfg, qCreatedUsers)
	}
	for i := 1; i <= numUsers; i++ {
		go o.activateUser(o.ctx, o.cfg, qCreatedUsers, qActivatedUsers, qFailedUsers)
	}
	for i := 1; i <= numUsers; i++ {
		go o.tokenUser(o.ctx, o.cfg, qActivatedUsers)
	}

	go func() {
		for {
			o.logger.Info().Msgf(
				"Concurrent queue len: | %6d | testUser creation queue:  %6d | testUser activation queue: %6d \n",
				len(concurrentGoroutines), len(qCreatedUsers), len(qActivatedUsers))
			if len(concurrentGoroutines) == 0 {
				done <- true
				o.logger.Info().Msg("Queues depleted, closing")
				break
			}
			time.Sleep(time.Second)
		}
	}()

	<-done
}

func (o *Ordered) createUser(ctx context.Context, cfg *config.Config, created chan testUser) {
	concurrentGoroutines <- struct{}{}

	pUsername, _ := generator.RandomString(usernameLen)
	email := pUsername + "@interview.com"

	newUser := &models.SignUpUserRequest{
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
	}
	marshaled, err := json.Marshal(newUser)
	if err != nil {
		o.logger.Info().Msgf("impossible to marshall: %+v\n", err)
	}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", o.endpointSignUp, bytes.NewReader(marshaled))
	if err != nil {
		o.logger.Error().Err(err).Msgf("impossible to create request: %s", err)
	}
	req.SetBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		errC := Body.Close()
		if errC != nil {
			o.logger.Error().Err(errC).Msg("ReadAll errC")
		}
	}(resp.Body)
	if err != nil {
		o.logger.Error().Err(err).Msg("Do err")
		<-concurrentGoroutines
		return
	}

	created <- testUser{
		Username: pUsername,
		Email:    email,
		Password: password,
	}
	<-concurrentGoroutines
}

func (o *Ordered) activateUser(ctx context.Context, cfg *config.Config, created chan testUser, activated chan testUser, failed chan testUser) {
	concurrentGoroutines <- struct{}{}
	user := <-created
	reqUser := &models.SignInUserRequest{Email: user.Email, Password: user.Password}

	marshaled, err := json.Marshal(reqUser)
	if err != nil {
		o.logger.Error().Err(err).Msgf("impossible to marshall: %s", err)
	}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", o.endpointAuthCode, bytes.NewReader(marshaled))
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		errC := Body.Close()
		if errC != nil {
			o.logger.Error().Err(errC).Msg("ReadAll errC")
		}
	}(resp.Body)
	if err != nil {
		o.logger.Error().Err(err).Msg("client Do")
	}
	if resp.StatusCode != http.StatusOK {
		o.logger.Info().Msgf("%s req body: %s\n", o.endpointAuthCode, string(marshaled))
		bodyBytes, errIo := io.ReadAll(resp.Body)
		if errIo != nil {
			o.logger.Error().Err(errIo).Msg("ReadAll errIo")
		}
		bodyString := string(bodyBytes)
		o.logger.Info().Msgf("%s body: %s", o.endpointAuthCode, bodyString)
	}

	type vCode struct {
		Token string `json:"verification_token"`
	}
	type target struct {
		User vCode `json:"user"`
	}
	tgt := target{}
	err = json.NewDecoder(resp.Body).Decode(&tgt)
	if err != nil {
		o.logger.Error().Err(err).Msg("Decoder err")
		<-concurrentGoroutines
		return
	}
	defer func(Body io.ReadCloser) {
		errC := Body.Close()
		if errC != nil {
			o.logger.Error().Err(errC).Msg("ReadAll errC")
		}
	}(resp.Body)

	client = &http.Client{}
	req, err = http.NewRequestWithContext(
		ctx,
		"GET",
		o.endpointVerifyCode+tgt.User.Token,
		bytes.NewReader(marshaled),
	)
	if err != nil {
		o.logger.Error().Err(err).Msgf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)

	resp, err = client.Do(req)
	if err != nil {
		<-concurrentGoroutines
		o.logger.Error().Err(err).Msg("client Do")
		return
	}
	defer func(Body io.ReadCloser) {
		errIo := Body.Close()
		if errIo != nil {
			o.logger.Error().Err(err).Msg("close body")
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		o.logger.Info().Msgf("%s req body: %s\n", o.endpointVerifyCode, string(marshaled))
		bodyBytes, errIo := io.ReadAll(resp.Body)
		if errIo != nil {
			o.logger.Error().Err(errIo).Msg("ReadAll err")
		}
		bodyString := string(bodyBytes)
		o.logger.Info().Msgf("%s body: %s", o.endpointVerifyCode, bodyString)
	}
	if err != nil {
		o.logger.Error().Err(err).Msg("verify err")
		<-concurrentGoroutines
		return
	}
	defer func(Body io.ReadCloser) {
		errC := Body.Close()
		if errC != nil {
			o.logger.Error().Err(errC).Msg("ReadAll errC")
		}
	}(resp.Body)
	activated <- testUser{
		ValidationCode: tgt.User.Token,
		Username:       user.Username,
		Email:          user.Email,
		Password:       user.Password,
	}
	<-concurrentGoroutines
}

func (o *Ordered) tokenUser(ctx context.Context, cfg *config.Config, activated chan testUser) {
	concurrentGoroutines <- struct{}{}

	user := <-activated
	credentials := &models.SignInUserRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	marshaled, err := json.Marshal(credentials)
	if err != nil {
		o.logger.Error().Err(err).Msgf("impossible to marshall: %s", err)
	}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", o.endpointSignIn, bytes.NewReader(marshaled))
	if err != nil {
		o.logger.Error().Err(err).Msgf("req creation err: %s", err)
	}

	req.SetBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	if err != nil {
		o.logger.Error().Err(err).Msg("Request Err")
	}
	defer func(Body io.ReadCloser) {
		errC := Body.Close()
		if errC != nil {
			o.logger.Error().Err(errC).Msg("ReadAll errC")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		o.logger.Info().Msgf("req body: %s\n", string(marshaled))
		o.logger.Info().Msgf("resp: %#v\n body: %s", resp, resp.Body)
		bodyBytes, errIo := io.ReadAll(resp.Body)
		if errIo != nil {
			o.logger.Error().Err(errIo).Msg("read body")
		}
		bodyString := string(bodyBytes)
		o.logger.Info().Msg("body: " + bodyString)
	}

	defer func(Body io.ReadCloser) {
		errIo := Body.Close()
		if errIo != nil {
			o.logger.Error().Err(errIo).Msg("close body")
		}
	}(resp.Body)
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		o.logger.Info().Msg("ReadAll err: " + err.Error())
	}
	<-concurrentGoroutines
}
