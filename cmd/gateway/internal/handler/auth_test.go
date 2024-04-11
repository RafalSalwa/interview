//go:build integration
// +build integration

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/config"
	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/internal/cqrs"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	"github.com/RafalSalwa/interview-app-srv/pkg/responses"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	password = "VeryG00dPass!"
)

var (
	cfg     *config.Config
	service *cqrs.Application
	l       *logger.Logger

	handler AuthHandler

	email       string
	vCode       string
	accessToken string
	reqSignUp   models.SignUpUserRequest
	reqSignIn   models.SignInUserRequest

	once sync.Once
)

func initAuthHandler(t *testing.T) {

	once.Do(func() {
		rand.New(rand.NewSource(time.Now().UnixNano()))
		var err error
		cfg, err = config.InitConfig()
		assert.NoError(t, err)

		service, err = cqrs.NewService(cfg.Grpc)
		assert.NoError(t, err)

		l = logger.NewConsole()
		handler = NewAuthHandler(service, l)

		email = fmt.Sprintf("interview%d@interview.com", rand.Int31())
		reqSignUp = models.SignUpUserRequest{
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
		}

		reqSignIn = models.SignInUserRequest{
			Email:    email,
			Password: password,
		}
	})
}

func TestSignUpUser(t *testing.T) {

	initAuthHandler(t)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqSignUp)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/auth/signup", &buf)
	w := httptest.NewRecorder()

	http.HandlerFunc(handler.SignUpUser()).ServeHTTP(w, r)
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := `{"status":"created"}`
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
	}
}
func TestGetVerificationCode(t *testing.T) {
	initAuthHandler(t)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqSignIn)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/auth/code", &buf)
	w := httptest.NewRecorder()

	http.HandlerFunc(handler.GetVerificationCode()).ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var jsonResp responses.UserResponse
	err = json.Unmarshal([]byte(w.Body.String()), &jsonResp)
	assert.NoError(t, err)
	vCode = jsonResp.VerificationCode
}

func TestVerify(t *testing.T) {
	initAuthHandler(t)
	r := httptest.NewRequest(http.MethodGet, "/verify?code="+vCode, nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(handler.Verify()).ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status":"ok"}`
	if w.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
	}
}

func TestSignInUser(t *testing.T) {
	initAuthHandler(t)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqSignIn)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/auth/signin", &buf)
	w := httptest.NewRecorder()
	http.HandlerFunc(handler.SignInUser()).ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if !strings.Contains(w.Body.String(), `"token"`) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
	}

	var jsonResp responses.UserResponse
	err = json.Unmarshal([]byte(w.Body.String()), &jsonResp)
	assert.NoError(t, err)
	accessToken = jsonResp.AccessToken
}
