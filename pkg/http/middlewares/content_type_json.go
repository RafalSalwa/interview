package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ContentTypeJSON() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// if r.Header.Get("Content-type") != "application/json" {
			//	w.WriteHeader(http.StatusUnsupportedMediaType)
			//	_, err := w.Write([]byte("415 - Unsupported Media Type. Only JSON files are allowed"))
			//	if err != nil {
			//		return
			//	}
			//	return
			//}

			w.Header().Set("Content-Type", "application/json;charset=utf8")

			h.ServeHTTP(w, r)
		})
	}
}
