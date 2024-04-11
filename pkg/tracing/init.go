package tracing

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer creates a new tracer for a service.
func InitTracer(serviceName, instanceName string) (trace.Tracer, func(), error) {
	tp, err := NewOTELTracerProvider(serviceName, instanceName)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't initialize tracer provider: %w", err)
	}

	// Cleanly shutdown and flush telemetry when the application exits.
	shutdown := func() {
		// Do not make the application hang when it is shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err = tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}

	return tp.Tracer(serviceName), shutdown, err
}

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func NewOTELTracerProvider(serviceName, instanceName string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://otel_collector:14278/api/traces")))
	if err != nil {
		return nil, err
	}
	return tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("service-instance", instanceName),
		)),
	), nil
}
