package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	SuccessHTTPRequests prometheus.Counter
	ErrorHTTPRequests   prometheus.Counter
}

func NewApiGatewayMetrics(serviceName string) *Metrics {
	return &Metrics{
		SuccessHTTPRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", serviceName),
			Help: "The total number of success http requests",
		}),
		ErrorHTTPRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", serviceName),
			Help: "The total number of error http requests",
		}),
	}
}
