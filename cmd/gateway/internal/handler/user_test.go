//go:build integration
// +build integration

package handler

import (
	"encoding/json"
	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/config"
	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/internal/cqrs"
	"github.com/RafalSalwa/interview-app-srv/pkg/jwt"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/RafalSalwa/interview-app-srv/pkg/responses"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

var uHandler UserHandler

type userTest struct {
	cfg      *config.Config
	service  *cqrs.Application
	l        *logger.Logger
	aHandler AuthHandler
	uHandler UserHandler

	accessToken string

	once sync.Once
}

func initUserHandler(t *testing.T) *userTest {
	ut := &userTest{}
	ut.once.Do(func() {
		var err error
		cfg, err := config.InitConfig()
		assert.NoError(t, err)
		ut.cfg = cfg

		service, err := cqrs.NewService(cfg.Grpc)
		assert.NoError(t, err)
		ut.service = service
		l := logger.NewConsole()
		ut.l = l
		ut.aHandler = NewAuthHandler(service, l)
		ut.uHandler = NewUserHandler(service, l)

	})
	return ut
}

func TestGetUserById(t *testing.T) {

	ut := initUserHandler(t)
	r := httptest.NewRequest(http.MethodPost, "/user", nil)
	w := httptest.NewRecorder()

	sub, err := jwt.ValidateToken(accessToken, ut.cfg.Auth.JWTToken.Access.PublicKey)
	assert.NoError(t, err)
	r.Header.Set("x-user-id", strconv.FormatInt(sub.ID, 10))

	http.HandlerFunc(ut.uHandler.GetUserByID()).ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var jsonResp responses.UserResponse
	err = json.Unmarshal([]byte(w.Body.String()), &jsonResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonResp)
	assert.NotEmpty(t, jsonResp.Id)
	assert.NotEmpty(t, jsonResp.Username)
	assert.Equal(t, true, jsonResp.Verified)
}
