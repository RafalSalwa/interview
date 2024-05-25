package handler

import (
	"net/http"

	"github.com/RafalSalwa/auth-api/pkg/encdec"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/responses"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type (
	CryptonHandler interface {
		RouteRegisterer

		Encrypt() http.HandlerFunc
		Decrypt() http.HandlerFunc
	}
	cryptonHandler struct {
		logger *logger.Logger
	}
)

func NewCryptonHandler(l *logger.Logger) CryptonHandler {
	return cryptonHandler{l}
}

func (h cryptonHandler) RegisterRoutes(r *mux.Router, cfg interface{}) {
	r.Methods(http.MethodGet).Path("/encrypt/{message}").HandlerFunc(h.Encrypt())
	r.Methods(http.MethodGet).Path("/decrypt/{message}").HandlerFunc(h.Decrypt())
}

func (h cryptonHandler) Encrypt() http.HandlerFunc {
	var message string
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.GetTracerProvider().Tracer("Encrypt").Start(r.Context(), "Encrypt Handler")
		defer span.End()
		message = mux.Vars(r)["message"]

		encrypted := encdec.Encrypt(message)
		responses.RespondString(w, encrypted)
	}
}

func (h cryptonHandler) Decrypt() http.HandlerFunc {
	var message string
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.GetTracerProvider().Tracer("Encrypt").Start(r.Context(), "Encrypt Handler")
		defer span.End()
		message = mux.Vars(r)["message"]

		encrypted, _ := encdec.Decrypt(message)
		responses.RespondString(w, encrypted)
	}
}
