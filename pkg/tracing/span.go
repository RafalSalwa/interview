package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func RecordError(span trace.Span, err error) {
	span.RecordError(err, trace.WithStackTrace(true))
	span.SetStatus(codes.Error, err.Error())
}

func InitSpan(ctx context.Context, tracer, spanName string) (context.Context, trace.Span) {
	return otel.GetTracerProvider().Tracer(tracer).Start(ctx, spanName)
}
