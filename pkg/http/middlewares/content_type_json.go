package middlewares

import (
	"net/http"

	"github.com/RafalSalwa/auth-api/pkg/responses"
	"github.com/gorilla/mux"
)

func ContentTypeJSON() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-type") != "application/json" &&
				r.Method == http.MethodPost {
				responses.RespondInternalServerError(w)
				return
			}

			w.Header().Set("Content-Type", "application/json;charset=utf8")

			h.ServeHTTP(w, r)
		})
	}
}
