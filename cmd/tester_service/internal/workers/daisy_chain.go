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
	cfg    *config.Config
	logger *logger.Logger
}

const numChannels = 4

func NewDaisyChain(cfg *config.Config, l *logger.Logger) WorkerRunner {
	return &DaisyChain{
		cfg:    cfg,
		logger: l,
	}
}

func (s DaisyChain) Run() {
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

func (s DaisyChain) dcCreateUser(ctx context.Context) testUser {
	pUsername, _ := generator.RandomString(12)
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
		s.logger.Error().Err(err).Msgf("impossible to marshall: %s", err)
	}
	client := &http.Client{}
	URL := fmt.Sprintf("http://%s/auth/signup", s.cfg.HTTP.Addr)
	// pass the values to the request's body
	req, err := http.NewRequestWithContext(ctx, "POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		s.logger.Error().Err(err).Msgf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Msg("Do err")
	}
	defer func(Body io.ReadCloser) {
		errC := Body.Close()
		if errC != nil {
			s.logger.Error().Err(errC).Msg("ReadAll errC")
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		s.logger.Error().Msgf("    %s req body: %s\n", URL, string(marshaled))
		bodyBytes, errIo := io.ReadAll(resp.Body)
		if errIo != nil {
			s.logger.Error().Err(errIo).Msgf("impossible to marshall: %s\n", errIo)
		}
		bodyString := string(bodyBytes)
		s.logger.Info().Msgf("    %s body: %s", URL, bodyString)
	} else {
		s.logger.Info().Msgf(color.GreenString("OK"))
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
