package auth

import (
	"net/http"
	"strings"

	"github.com/RafalSalwa/interview-app-srv/pkg/responses"
)

type bearerTokenHandler struct {
	token string
}

func (a *bearerTokenHandler) Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer "
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)
		if authHeader == "" || reqToken == authHeader {
			responses.RespondNotAuthorized(w, "Authentication header not present or malformed")
			return
		}
		h(w, r)
	}
}

func newBearerTokenMiddleware(token string) *bearerTokenHandler {
	return &bearerTokenHandler{
		token: token,
	}
}
