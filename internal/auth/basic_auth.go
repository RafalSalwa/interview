package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"os"

	apiHandler "github.com/RafalSalwa/interview-app-srv/api/handler"
	"github.com/RafalSalwa/interview-app-srv/api/resource/responses"
)

type basicAuth struct {
	h        apiHandler.HandlerFunc
	username string
	password string
}

func (a *basicAuth) middleware(h apiHandler.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		authUsername := os.Getenv("AUTH_USERNAME")
		authPassword := os.Getenv("AUTH_PASSWORD")

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
		responses.RespondNotAuthorized(w, "")
		return
	}
}

func newBasicAuthMiddleware(h apiHandler.HandlerFunc, username string, password string) *basicAuth {
	return &basicAuth{
		h:        h,
		username: username,
		password: password,
	}
}
