package router

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/RafalSalwa/auth-api/docs"
	"github.com/RafalSalwa/auth-api/pkg/http/middlewares"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewHTTPRouter(l *logger.Logger) *mux.Router {
	router := mux.NewRouter()

	promMiddleware := middlewares.NewPrometheusMiddleware()

	router.Use(
		middlewares.ContentTypeJSON(),
		middlewares.CorrelationID(),
		middlewares.CORS(),
		middlewares.RequestLog(l),
		promMiddleware.Prometheus(),
	)
	router.Handle("/metrics", promhttp.Handler())

	setupSwagger(router)

	return router
}

func setupSwagger(r *mux.Router) {
	docs.SwaggerInfo.Title = "Interview API for Gateway Service"
	docs.SwaggerInfo.Description = "API Gateway that works like a backends for frontends pattern" +
		" and passes requests to specific services"

	jsonFile, err := os.Open("docs/swagger.json")
	if err != nil {
		log.Fatal(err)
	}

	bytesJSON, _ := io.ReadAll(jsonFile)
	docs.SwaggerInfo.SwaggerTemplate = string(bytesJSON)

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	r.Path("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusSeeOther)
	})).Methods(http.MethodGet)
}
