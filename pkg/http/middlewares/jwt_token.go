package middlewares

import (
	"net/http"
	"strings"

	"github.com/RafalSalwa/auth-api/pkg/http/auth"
	"github.com/RafalSalwa/auth-api/pkg/responses"

	"github.com/RafalSalwa/auth-api/pkg/jwt"
	"github.com/gorilla/mux"
)

func ValidateJWTAccessToken(c auth.JWTConfig) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			if authorizationHeader != "" {
				bearerToken := strings.Split(authorizationHeader, " ")

				if len(bearerToken) == 2 {
					sub, err := jwt.ValidateToken(bearerToken[1], c.Access.PublicKey)
					if err != nil {
						responses.RespondNotAuthorized(w, "Wrong access token")
						return
					}
					r.Header.Set("x-user-id", sub)
				}
			} else {
				responses.RespondNotAuthorized(w, "An authorization header is required")
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
