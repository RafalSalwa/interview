package handler

import (
	"net/http"

	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/gorilla/mux"
)

type DocsHandler interface {
	RouteRegisterer
	Documentation() HandlerFunc
}

type docsHandler struct {
	logger *logger.Logger
}

func (d docsHandler) Documentation() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusTemporaryRedirect)
	}
}

func (d docsHandler) RegisterRoutes(r *mux.Router, cfg interface{}) {
	r.Methods(http.MethodGet, http.MethodPost).Path("/").HandlerFunc(d.Documentation())
}

func NewDocsHandler(l *logger.Logger) DocsHandler {
	return docsHandler{l}
}
