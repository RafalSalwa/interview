package tracing

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetricsEndpoint(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		server := &http.Server{
			Addr:         addr,
			Handler:      nil,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		}

		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Failed to serve Prometheus metrics:", err)
		}
	}()
}
