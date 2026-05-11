// Package telemetry initializes OpenTelemetry tracing with the GCP Cloud Trace exporter.
package telemetry

import (
	"context"
	"fmt"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcpdetector "go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Init sets up the OpenTelemetry tracer provider with the GCP Cloud Trace
// exporter and GCP resource detection. It returns a shutdown function that
// should be deferred in main().
//
// The sampling ratio is configured via OTEL_TRACES_SAMPLER_ARG env var (default: 0.1).
// No OTEL_EXPORTER_OTLP_ENDPOINT is needed — the exporter calls the Cloud Trace API directly.
func Init(ctx context.Context) (func(context.Context) error, error) {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "swim-gen-backend"
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	opts := []texporter.Option{}
	if projectID != "" {
		opts = append(opts, texporter.WithProjectID(projectID))
	}

	exporter, err := texporter.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("creating Cloud Trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(gcpdetector.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Parse sampler arg (default 10% sampling)
	samplerArg := 0.1
	if s := os.Getenv("OTEL_TRACES_SAMPLER_ARG"); s != "" {
		if _, err := fmt.Sscanf(s, "%f", &samplerArg); err != nil {
			return nil, fmt.Errorf("parsing OTEL_TRACES_SAMPLER_ARG: %w", err)
		}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplerArg))),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

// TraceIDFromContext extracts the trace ID from the current span context.
// Returns an empty string if no active span exists or the trace ID is invalid.
func TraceIDFromContext(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if sc.HasTraceID() {
		return sc.TraceID().String()
	}
	return ""
}
