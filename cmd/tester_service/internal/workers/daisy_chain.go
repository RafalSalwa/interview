package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RafalSalwa/auth-api/pkg/logger"

	"github.com/RafalSalwa/auth-api/cmd/tester_service/config"
	"github.com/RafalSalwa/auth-api/pkg/generator"
	"github.com/RafalSalwa/auth-api/pkg/models"
	pb "github.com/RafalSalwa/auth-api/proto/grpc"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DaisyChain struct {
	cfg                *config.Config
	logger             *logger.Logger
	endpoint           string
	endpointSignUp     string
	endpointSignIn     string
	endpointAuthCode   string
	endpointVerifyCode string
}

const (
	numChannels = 4
)

func NewDaisyChain(cfg *config.Config, l *logger.Logger) WorkerRunner {
	dc := &DaisyChain{
		cfg:    cfg,
		logger: l,
	}
	dc.endpoint = fmt.Sprintf("http://%s", cfg.HTTP.Addr)
	dc.endpointSignUp = fmt.Sprintf("%s/auth/signup", dc.endpoint)
	dc.endpointSignIn = fmt.Sprintf("%s/auth/signin", dc.endpoint)
	dc.endpointAuthCode = fmt.Sprintf("%s/auth/code", dc.endpoint)
	dc.endpointVerifyCode = fmt.Sprintf("%s/auth/verify", dc.endpoint)

	return dc
}

func (s *DaisyChain) Run() {
	tasks := [numChannels]string{"signUp", "getCode", "activate", "signIn"}
	ctx := context.Background()

	leftmost := make(chan testUser)
	right := leftmost
	left := leftmost

	for i := 0; i < numChannels; i++ {
		go worker(ctx, left, right, tasks[i])
	}

	leftmost <- s.dcCreateUser(ctx)
}
func worker(ctx context.Context, in <-chan testUser, out chan<- testUser, task string) {
	inUser := <-in

	var outUser testUser
	switch task {
	case "activate":
		outUser = dcActivateUser(ctx, inUser)
	case "token":
		outUser = dcTokenUser(ctx, inUser)
	}
	out <- outUser
}

func (dc *DaisyChain) dcCreateUser(ctx context.Context) testUser {
	pUsername, _ := generator.RandomString(usernameLen)
	email := pUsername + emailDomain

	user := testUser{
		Username: pUsername,
		Email:    email,
		Password: password,
	}

	newUser := &models.SignUpUserRequest{
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.Password,
	}
	marshaled, err := json.Marshal(newUser)
	if err != nil {
		dc.logger.Error().Err(err).Msgf("impossible to marshall: %s", err)
	}
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", dc.endpointSignUp, bytes.NewReader(marshaled))
	if err != nil {
		dc.logger.Error().Err(err).Msgf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(dc.cfg.Auth.BasicAuth.Username, dc.cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		errC := Body.Close()
		if errC != nil {
			dc.logger.Error().Err(errC).Msg("ReadAll errC")
		}
	}(resp.Body)
	if err != nil {
		dc.logger.Error().Err(err).Msg("Do err")
	}

	if resp.StatusCode != http.StatusCreated {
		dc.logger.Error().Msgf("    %s req body: %s\n", dc.endpointSignUp, string(marshaled))
		bodyBytes, errIo := io.ReadAll(resp.Body)
		if errIo != nil {
			dc.logger.Error().Err(errIo).Msgf("impossible to marshall: %s\n", errIo)
		}
		bodyString := string(bodyBytes)
		dc.logger.Info().Msgf("    %s body: %s", dc.endpointSignUp, bodyString)
	} else {
		dc.logger.Info().Msgf(color.GreenString("OK"))
	}
	return user
}

func dcActivateUser(ctx context.Context, inUser testUser) testUser {
	rVerification := &pb.VerifyUserRequest{Code: inUser.ValidationCode}
	conn, _ := grpc.Dial("0.0.0.0:8022", grpc.WithTransportCredentials(insecure.NewCredentials()))

	userClient := pb.NewUserServiceClient(conn)

	_, _ = userClient.VerifyUser(ctx, rVerification)
	return inUser
}

func dcTokenUser(ctx context.Context, inUser testUser) testUser {
	conn, _ := grpc.Dial("0.0.0.0:8032", grpc.WithTransportCredentials(insecure.NewCredentials()))

	authClient := pb.NewAuthServiceClient(conn)

	credentials := &pb.SignInUserInput{
		Username: inUser.Username,
		Password: inUser.Password,
	}
	_, _ = authClient.SignInUser(ctx, credentials)
	return inUser
}
