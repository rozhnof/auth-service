package tracer

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewJaegerExporter(ctx context.Context, url string) (tracesdk.SpanExporter, error) {
	return otlptracehttp.New(
		ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(url),
	)
}
