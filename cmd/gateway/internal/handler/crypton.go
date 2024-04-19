package handler

import (
	"net/http"

	"github.com/RafalSalwa/interview-app-srv/pkg/encdec"
	"github.com/RafalSalwa/interview-app-srv/pkg/responses"
	"go.opentelemetry.io/otel"

	"github.com/gorilla/mux"

	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
)

type CryptonHandler interface {
	RouteRegisterer

	Encrypt() http.HandlerFunc
	Decrypt() http.HandlerFunc
}

type cryptonHandler struct {
	logger *logger.Logger
}

func (c cryptonHandler) RegisterRoutes(r *mux.Router, cfg interface{}) {
	r.Methods(http.MethodGet).Path("/encrypt/{message}").HandlerFunc(c.Encrypt())
	r.Methods(http.MethodGet).Path("/decrypt/{message}").HandlerFunc(c.Decrypt())
}

func (c cryptonHandler) Encrypt() http.HandlerFunc {
	var message string
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.GetTracerProvider().Tracer("Encrypt").Start(r.Context(), "Encrypt Handler")
		defer span.End()
		message = mux.Vars(r)["message"]

		encrypted := encdec.Encrypt(message)
		responses.RespondString(w, encrypted)
	}
}

func (c cryptonHandler) Decrypt() http.HandlerFunc {
	var message string
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.GetTracerProvider().Tracer("Encrypt").Start(r.Context(), "Encrypt Handler")
		defer span.End()
		message = mux.Vars(r)["message"]

		encrypted, _ := encdec.Decrypt(message)
		responses.RespondString(w, encrypted)
	}
}

func NewCryptonHandler(l *logger.Logger) CryptonHandler {
	return cryptonHandler{l}
}
