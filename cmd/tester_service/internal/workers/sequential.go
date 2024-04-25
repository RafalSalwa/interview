package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RafalSalwa/auth-api/pkg/logger"

	"github.com/RafalSalwa/auth-api/cmd/tester_service/config"
	"github.com/RafalSalwa/auth-api/pkg/generator"
	"github.com/RafalSalwa/auth-api/pkg/models"
	"github.com/fatih/color"
)

type Sequential struct {
	ctx    context.Context
	logger *logger.Logger
	cfg    *config.Config
	client *http.Client
}

const (
	usernameLen = 12
)

func NewSequential(ctx context.Context, cfg *config.Config, l *logger.Logger) WorkerRunner {
	return &Sequential{
		ctx:    ctx,
		logger: l,
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (s Sequential) Run() {
	for {
		pUsername, _ := generator.RandomString(usernameLen)
		email := pUsername + emailDomain

		user := &testUser{
			Username: pUsername,
			Email:    email,
			Password: password,
		}
		s.signUp(user)
		s.getVerificationCode(user)
		s.activateUser(user)
		s.signIn(user)

		time.Sleep(10 * time.Second)
	}
}

func (s Sequential) signUp(user *testUser) {
	newUser := &models.SignUpUserRequest{
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.Password,
	}

	marshaled, err := json.Marshal(newUser)
	if err != nil {
		s.logger.Error().Err(err).Msg("impossible to marshall")
	}

	client := &http.Client{}
	URL := fmt.Sprintf("http://%s/auth/signup", s.cfg.HTTP.Addr)
	req, err := http.NewRequestWithContext(s.ctx, "POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		s.logger.Error().Err(err).Msgf("impossible to read all body of response: %s", err)
	}

	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Msgf("req do err: ", err)
	}
	if resp.StatusCode != http.StatusCreated {
		s.logger.Error().Err(err).Msgf("    %s req body: %s\n", URL, string(marshaled))
		bodyBytes, errIo := io.ReadAll(resp.Body)
		if errIo != nil {
			s.logger.Error().Err(errIo).Msgf("impossible to marshall: %s\n", err)
		}
		bodyString := string(bodyBytes)
		s.logger.Error().Err(err).Msgf("    %s body: %s", URL, bodyString)
	} else {
		s.logger.Info().Msgf(color.GreenString("OK"))
	}
	err = resp.Body.Close()
	if err != nil {
		s.logger.Error().Err(err).Msgf("req do err: ", err)
	}
}

func (s Sequential) getVerificationCode(user *testUser) {
	reqUser := &models.SignInUserRequest{Email: user.Email, Password: user.Password}

	marshaled, err := json.Marshal(reqUser)
	if err != nil {
		s.logger.Error().Err(err).Msgf("impossible to marshall: %s\n", err)
	}
	URL := "http://" + s.cfg.HTTP.Addr + "/auth/code"
	req, err := http.NewRequestWithContext(s.ctx, "POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		s.logger.Error().Err(err).Msgf("impossible to read all body of response: %s\n", err)
	}
	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Msgf("impossible to marshall: %s\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		s.logger.Error().Msgf("%s req body: %s\n", URL, string(marshaled))
		bodyBytes, errIo := io.ReadAll(resp.Body)
		if errIo != nil {
			s.logger.Error().Err(errIo).Msgf("Incorrect response: %s", errIo)
		}
		bodyString := string(bodyBytes)
		s.logger.Error().Msgf("%s body: %s", URL, bodyString)
	} else {
		s.logger.Info().Msgf(color.GreenString("OK "))
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
		s.logger.Error().Err(err).Msg("impossible to unmarshall")
		return
	}
	user.ValidationCode = tgt.User.Token
	err = resp.Body.Close()
	if err != nil {
		s.logger.Error().Err(err).Msg("impossible to unmarshall")
		return
	}
}

func (s Sequential) activateUser(user *testUser) {
	URL := "http://" + s.cfg.HTTP.Addr + "/auth/verify/" + user.ValidationCode
	req, err := http.NewRequestWithContext(s.ctx, "GET", URL, nil)
	if err != nil {
		s.logger.Error().Err(err).Msgf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Msgf("/auth/verify/", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Error().Msgf("%s req :\n", URL)
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		s.logger.Error().Msgf("%s body: %s", URL, bodyString)
	} else {
		s.logger.Info().Msgf(color.GreenString("OK "))
	}
	err = resp.Body.Close()
	if err != nil {
		s.logger.Error().Err(err).Msg("Cannot close response body")
		return
	}
}

func (s Sequential) signIn(user *testUser) {
	credentials := &models.SignInUserRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	marshaled, err := json.Marshal(credentials)
	if err != nil {
		s.logger.Error().Err(err).Msgf("impossible to marshall: %s\n", err)
	}
	URL := "http://" + s.cfg.HTTP.Addr + "/auth/signin"
	req, err := http.NewRequestWithContext(s.ctx, "POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		s.logger.Error().Err(err).Msgf("impossible to read all body of response: %s", err)
	}

	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Msgf("Do err: ", err)
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Error().Msgf("%s req body: %s\n", URL, string(marshaled))
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		s.logger.Error().Msgf("%s body: %s", URL, bodyString)
	} else {
		s.logger.Info().Msgf(color.GreenString("OK "))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error().Err(err).Msgf("ReadAll err: ", err)
	}
	err = resp.Body.Close()
	if err != nil {
		s.logger.Error().Err(err).Msg("Cannot close response body")
		return
	}
	s.logger.Info().Msgf("Token: ", string(respBody))
}
