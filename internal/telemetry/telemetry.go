package telemetry

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.opentelemetry.io/otel/trace"
)

// InitTracer sets up the global tracer with platform-specific attributes
func InitTracer() (*sdktrace.TracerProvider, error) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("terraform-provider-hello"),
			// Custom attributes for API convergence tracking
			attribute.String("platform.os", runtime.GOOS),
			attribute.String("platform.arch", runtime.GOARCH),
			attribute.String("sdk.version", "1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}
	// 1. Define the destination (Exporter)
	// By default, this looks for a collector at localhost:4317
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		// In production, use an exporter like OTLP or Jaeger
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("hello-provider").Start(ctx, name)
}
