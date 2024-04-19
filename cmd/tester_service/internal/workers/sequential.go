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

	"github.com/RafalSalwa/interview-app-srv/cmd/tester_service/config"
	"github.com/RafalSalwa/interview-app-srv/pkg/generator"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	"github.com/fatih/color"
)

type Sequential struct {
	ctx    context.Context
	cfg    *config.Config
	client *http.Client
}

func NewSequential(ctx context.Context, cfg *config.Config) WorkerRunner {
	return &Sequential{
		ctx:    ctx,
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (s Sequential) Run() {
	for {
		pUsername, _ := generator.RandomString(12)
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
		fmt.Println()
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
		log.Printf("impossible to marshall: %+v\n", err)
	}

	client := &http.Client{}
	URL := fmt.Sprintf("http://%s/auth/signup", s.cfg.HTTP.Addr)
	req, err := http.NewRequest("POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		fmt.Printf("impossible to read all body of response: %s", err)
	}

	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("req do err: ", err)
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
	err = resp.Body.Close()
	if err != nil {
		fmt.Println("req do err: ", err)
	}
}

func (s Sequential) getVerificationCode(user *testUser) {
	reqUser := &models.SignInUserRequest{Email: user.Email, Password: user.Password}

	marshaled, err := json.Marshal(reqUser)
	if err != nil {
		fmt.Printf("impossible to marshall: %s\n", err)
	}
	URL := "http://" + s.cfg.HTTP.Addr + "/auth/code"
	req, err := http.NewRequestWithContext(s.ctx, "POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		fmt.Printf("impossible to read all body of response: %s\n", err)
	}
	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)
	fmt.Print("Fetching vCode: ")
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Printf("impossible to marshall: %s\n", err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Err")
		fmt.Printf("%s req body: %s\n", URL, string(marshaled))
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("impossible to marshall: %s\n", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("%s body: %s", URL, bodyString)
	} else {
		fmt.Print(color.GreenString("OK "))
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
		fmt.Println(err)
		return
	}
	user.ValidationCode = tgt.User.Token
	fmt.Println(tgt.User.Token)
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s Sequential) activateUser(user *testUser) {
	URL := "http://" + s.cfg.HTTP.Addr + "/auth/verify/" + user.ValidationCode
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Printf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)

	fmt.Print("Veryfing: ")
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println("/auth/verify/", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("Err ")
		fmt.Printf("%s req :\n", URL)
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("%s body: %s", URL, bodyString)
	} else {
		fmt.Println(color.GreenString("OK "))
	}
	if err != nil {
		fmt.Println("verify err:", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
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
		fmt.Printf("impossible to marshall: %s\n", err)
	}
	URL := "http://" + s.cfg.HTTP.Addr + "/auth/signin"
	req, err := http.NewRequest("POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		fmt.Printf("impossible to read all body of response: %s", err)
	}

	req.SetBasicAuth(s.cfg.Auth.BasicAuth.Username, s.cfg.Auth.BasicAuth.Password)
	fmt.Print("Logging In: ")
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println("Do err: ", err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Err ")
		fmt.Printf("%s req body: %s\n", URL, string(marshaled))
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("%s body: %s", URL, bodyString)
	} else {
		fmt.Println(color.GreenString("OK "))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll err: ", err)
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Token: ", string(respBody))
}
