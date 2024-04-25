package main

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
	"github.com/RafalSalwa/auth-api/cmd/tester_service/internal/workers"
	"github.com/RafalSalwa/auth-api/pkg/generator"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/models"
)

const (
	emailDomain               = "@interview.com"
	password                  = "VeryG00dPass!"
	numUsers                  = 20
	maxNbConcurrentGoroutines = 10
)

type testUser struct {
	ValidationCode string
	Username       string
	Email          string
	Password       string
}

var (
	concurrentGoroutines = make(chan struct{}, maxNbConcurrentGoroutines)
)

func main() {
	worker := workers.NewWorker("daisy_chain")
	worker.Run()
}

func runWorkersInOrder(ctx context.Context, cfg *config.Config, l *logger.Logger) {
	qCreatedUsers := make(chan testUser, numUsers)
	qActivatedUsers := make(chan testUser, numUsers)
	qFailedUsers := make(chan testUser, numUsers)

	done := make(chan bool)

	for i := 1; i <= numUsers; i++ {
		go createUser(ctx, cfg, qCreatedUsers, qFailedUsers)
	}
	for i := 1; i <= numUsers; i++ {
		go activateUser(ctx, cfg, qCreatedUsers, qActivatedUsers, qFailedUsers)
	}
	for i := 1; i <= numUsers; i++ {
		go tokenUser(ctx, cfg, qActivatedUsers, qFailedUsers)
	}

	go func() {
		for {
			fmt.Printf("Concurrent queue len: | %6d | testUser creation queue:  %6d | testUser activation queue: %6d \n", len(concurrentGoroutines), len(qCreatedUsers), len(qActivatedUsers))
			if len(concurrentGoroutines) == 0 {
				done <- true
				fmt.Println("Queues depleted, closing")
				break
			}
			time.Sleep(time.Second)
		}
	}()

	<-done
}

func createUser(ctx context.Context, cfg *config.Config, created, failed chan testUser) {
	concurrentGoroutines <- struct{}{}

	pUsername, _ := generator.RandomString(12)
	email := pUsername + emailDomain

	newUser := &models.SignUpUserRequest{
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
	}
	marshaled, err := json.Marshal(newUser)
	if err != nil {
		fmt.Printf("impossible to marshall: %+v\n", err)
	}
	client := &http.Client{}
	URL := fmt.Sprintf("http://%s/auth/signup", cfg.HTTP.Addr)
	req, err := http.NewRequestWithContext(ctx, "POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("req do err: ", err)
		<-concurrentGoroutines
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("body close err", err)
		}
	}(resp.Body)

	created <- testUser{
		Username: pUsername,
		Email:    email,
		Password: password,
	}
	<-concurrentGoroutines
}

func activateUser(ctx context.Context, cfg *config.Config, created chan testUser, activated chan testUser, failed chan testUser) {
	concurrentGoroutines <- struct{}{}
	user := <-created
	reqUser := &models.SignInUserRequest{Email: user.Email, Password: user.Password}

	marshaled, err := json.Marshal(reqUser)
	if err != nil {
		log.Fatalf("impossible to marshall: %s", err)
	}
	client := &http.Client{}
	URL := "http://" + cfg.HTTP.Addr + "/auth/code"
	req, err := http.NewRequestWithContext(ctx, "POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("%s req body: %s\n", URL, string(marshaled))
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("%s body: %s", URL, bodyString)
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
		<-concurrentGoroutines
		return
	}
	defer resp.Body.Close()
	if err != nil {
		<-concurrentGoroutines
		return
	}

	client = &http.Client{}
	URL = "http://" + cfg.HTTP.Addr + "/auth/verify/" + tgt.User.Token
	req, err = http.NewRequestWithContext(ctx, "GET", URL, bytes.NewReader(marshaled))
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	req.SetBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("/auth/verify/", err)
		<-concurrentGoroutines
		log.Fatal(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("defer close: ", err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Printf("%s req body: %s\n", URL, string(marshaled))
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("%s body: %s", URL, bodyString)
	}
	if err != nil {
		fmt.Println("verify err:", err)
		<-concurrentGoroutines
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	activated <- testUser{ValidationCode: tgt.User.Token, Username: user.Username, Email: user.Email, Password: user.Password}
	<-concurrentGoroutines
}

func tokenUser(ctx context.Context, cfg *config.Config, activated, failed chan testUser) {
	concurrentGoroutines <- struct{}{}

	user := <-activated
	credentials := &models.SignInUserRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	marshaled, err := json.Marshal(credentials)
	if err != nil {
		log.Fatalf("impossible to marshall: %s", err)
	}
	client := &http.Client{}
	URL := "http://" + cfg.HTTP.Addr + "/auth/signin"
	req, err := http.NewRequestWithContext(ctx, "POST", URL, bytes.NewReader(marshaled))
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}

	req.SetBasicAuth(cfg.Auth.BasicAuth.Username, cfg.Auth.BasicAuth.Password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Do err: ", err)
		log.Fatalf("impossible to read all body of response: %s", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("req body: %s\n", string(marshaled))
		fmt.Printf("resp: %#v\n body: %s", resp, resp.Body)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println("body: ", bodyString)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll err: ", err)
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	<-concurrentGoroutines
}
