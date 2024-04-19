package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/RafalSalwa/interview-app-srv/cmd/tester_service/config"
	"github.com/RafalSalwa/interview-app-srv/pkg/generator"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	pb "github.com/RafalSalwa/interview-app-srv/proto/grpc"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const numChannels = 4

type DaisyChain struct {
	cfg *config.Config
}

func NewDaisyChain(cfg *config.Config) WorkerRunner {
	return &DaisyChain{
		cfg: cfg,
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

	leftmost <- s.dcCreateUser(ctx, s.cfg)
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

func (s DaisyChain) dcCreateUser(ctx context.Context, cfg *config.Config) testUser {
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
		log.Fatalf("impossible to marshall: %s", err)
	}
	client := &http.Client{}
	URL := fmt.Sprintf("http://%s/auth/signup", s.cfg.HTTP.Addr)
	// pass the values to the request's body
	req, err := http.NewRequest("POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Err.")

		fmt.Printf("    %s req body: %s\n", URL, string(marshaled))
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("impossible to marshall: %s\n", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("    %s body: %s", URL, bodyString)
	} else {
		fmt.Println(color.GreenString("OK"))
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
