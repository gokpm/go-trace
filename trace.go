package trace

import (
	"context"
	"os"

	otrace "go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

type Config struct {
	Name        string
	Environment string
	URL         string
	Sampling    float64
}

func Setup(ctx context.Context, config *Config) (otrace.Tracer, error) {
	httpOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpointURL(config.URL),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	}
	exporter, err := otlptracehttp.New(ctx, httpOpts...)
	if err != nil {
		return nil, err
	}
	processor := trace.NewBatchSpanProcessor(exporter)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	base := resource.Default()
	newResource := resource.NewWithAttributes(
		base.SchemaURL(),
		semconv.ServiceName(config.Name),
		semconv.DeploymentEnvironmentName(config.Environment),
		semconv.HostName(hostname),
	)
	mergedResource, err := resource.Merge(base, newResource)
	if err != nil {
		return nil, err
	}
	sampler := trace.ParentBased(trace.TraceIDRatioBased(config.Sampling))
	providerOpts := []trace.TracerProviderOption{
		trace.WithBatcher(exporter),
		trace.WithResource(mergedResource),
		trace.WithSampler(sampler),
		trace.WithSpanProcessor(processor),
	}
	return trace.NewTracerProvider(providerOpts...).Tracer(config.Name), nil
}
