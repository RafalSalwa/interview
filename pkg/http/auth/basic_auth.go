package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/RafalSalwa/auth-api/pkg/responses"
)

type basicAuth struct {
	username string
	password string
}

func (a *basicAuth) Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		authUsername := a.username
		authPassword := a.password
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(authUsername))
			expectedPasswordHash := sha256.Sum256([]byte(authPassword))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				h(w, r)
				return
			}
		}
		responses.RespondNotAuthorized(w, "Missing or invalid credentials")
	}
}

func newBasicAuthMiddleware(username, password string) *basicAuth {
	return &basicAuth{
		username: username,
		password: password,
	}
}
