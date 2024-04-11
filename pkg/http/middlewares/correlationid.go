package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type CorrID struct {
	value string
}

func CorrelationID() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			id := r.Header.Get("X-Correlation-Id")
			if id == "" {
				newid := uuid.New()
				id = newid.String()
			}
			ctx = context.WithValue(ctx, CorrID{"correlation_id"}, id)
			r = r.WithContext(ctx)
			log := zerolog.Ctx(ctx)
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("correlation_id", id)
			})
			w.Header().Set("X-Correlation-Id", id)
			h.ServeHTTP(w, r)
		})
	}
}
