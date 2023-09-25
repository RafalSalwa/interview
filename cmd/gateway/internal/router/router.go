package router

import (
	_ "embed"
	"fmt"
	"github.com/RafalSalwa/interview-app-srv/docs"
	_ "github.com/RafalSalwa/interview-app-srv/docs"
	"github.com/RafalSalwa/interview-app-srv/pkg/http/middlewares"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"io"
	"net/http"
	"os"
)

func NewHTTPRouter(l *logger.Logger) *mux.Router {
	router := mux.NewRouter()

	promMiddleware := middlewares.NewPrometheusMiddleware()

	router.Use(
		middlewares.ContentTypeJson(),
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
	docs.SwaggerInfo.Description = "API Gateway that works like a backends for frontends pattern and passes requests to specific services"

	jsonFile, err := os.Open("docs/swagger.json")
	if err != nil {
		fmt.Println(err)
	}

	bytesJSON, _ := io.ReadAll(jsonFile)
	docs.SwaggerInfo.SwaggerTemplate = string(bytesJSON)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)
}
