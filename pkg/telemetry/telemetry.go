package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.28.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	serviceName    = "k8s-controller-tutorial"
	serviceVersion = "1.0.0"
)

// TracingConfig configuration for tracing
type TracingConfig struct {
	ServiceName    string
	ServiceVersion string
	EnableConsole  bool
}

// InitTracing initializes OpenTelemetry tracing
func InitTracing(ctx context.Context, config TracingConfig) (func(context.Context) error, error) {
	// Setup resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Setup trace exporter for stdout output
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Setup trace provider
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
	)

	// Set global trace provider
	otel.SetTracerProvider(tracerProvider)

	log.Info().Msg("OpenTelemetry tracing initialized")

	// Return shutdown function
	return tracerProvider.Shutdown, nil
}

// GetTracer returns a tracer for the given name
func GetTracer(name string) oteltrace.Tracer {
	return otel.Tracer(name)
}

// StartSpan creates a new span with logging
func StartSpan(ctx context.Context, tracer oteltrace.Tracer, spanName string, attrs ...attribute.KeyValue) (context.Context, oteltrace.Span) {
	ctx, span := tracer.Start(ctx, spanName)

	// Add attributes if provided
	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	// Log span start
	log.Info().
		Str("span_name", spanName).
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msg("Starting span")

	return ctx, span
}

// EndSpan ends a span with logging
func EndSpan(span oteltrace.Span) {
	defer span.End()

	log.Info().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msg("Ending span")
}

// AddSpanEvent adds an event to a span with logging
func AddSpanEvent(span oteltrace.Span, name string, attrs ...attribute.KeyValue) {
	span.AddEvent(name, oteltrace.WithAttributes(attrs...))

	log.Info().
		Str("span_name", name).
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msg("Added span event")
}

// SetSpanError sets an error for a span
func SetSpanError(span oteltrace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	log.Error().
		Err(err).
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msg("Span error recorded")
}

// RecordSpanDuration records operation duration
func RecordSpanDuration(span oteltrace.Span, operation string, duration time.Duration) {
	span.SetAttributes(
		attribute.String("operation", operation),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)

	log.Info().
		Str("operation", operation).
		Int64("duration_ms", duration.Milliseconds()).
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msg("Operation duration recorded")
}
