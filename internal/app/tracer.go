package app

import (
	"context"
	"fmt"

	tracing "github.com/rozhnof/auth-service/internal/infrastructure/tracer"
	"github.com/rozhnof/auth-service/internal/pkg/config"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func NewTracer(ctx context.Context, cfg config.Tracing, serviceName string) (trace.Tracer, func(context.Context) error, error) {
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	switch cfg.Output {
	case "stdout":
		exporter, err = tracing.NewConsoleExporter()
	case "jaeger":
		exporter, err = tracing.NewJaegerExporter(ctx, cfg.Endpoint)
	default:
		return nil, nil, fmt.Errorf("invalid tracing output parameter: %s", cfg.Output)
	}

	if err != nil {
		return nil, nil, err
	}

	provider, err := tracing.NewTraceProvider(exporter, serviceName)
	if err != nil {
		return nil, nil, err
	}

	otel.SetTracerProvider(provider)

	return provider.Tracer(cfg.Name), provider.Shutdown, nil
}
