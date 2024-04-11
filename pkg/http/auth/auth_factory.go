package auth

import (
	"errors"
	"net/http"
)

type IAuthType interface {
	Middleware(h http.HandlerFunc) http.HandlerFunc
}

type AuthType int

const (
	apiKey AuthType = 1 << iota
	basic
	bearerToken
)

var types = map[string]interface{}{
	"key":          apiKey,
	"basic":        basic,
	"bearer_token": bearerToken,
}

func NewAuthMethod(cfg Auth) (IAuthType, error) {
	val, ok := types[cfg.AuthMethod]
	if !ok {
		return nil, errors.New("wrong auth type")
	}
	switch val {
	case apiKey:
		return newAPIKeyMiddleware(cfg.APIKey), nil
	case basic:
		return newBasicAuthMiddleware(cfg.BasicAuth.Username, cfg.BasicAuth.Password), nil
	case bearerToken:
		return newBearerTokenMiddleware(cfg.BearerToken), nil
	}
	return nil, nil
}
