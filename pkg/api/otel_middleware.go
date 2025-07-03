package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JRaver/k8s-controller-tutorial/pkg/telemetry"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// OtelMiddleware creates middleware for OpenTelemetry tracing
func OtelMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	tracer := telemetry.GetTracer("k8s-controller-tutorial/api")

	return func(ctx *fasthttp.RequestCtx) {
		// Create context for span
		reqCtx := context.Background()

		// Create span for HTTP request
		spanName := fmt.Sprintf("%s %s", string(ctx.Method()), string(ctx.Path()))

		// Start span with attributes
		reqCtx, span := telemetry.StartSpan(reqCtx, tracer, spanName,
			attribute.String("http.method", string(ctx.Method())),
			attribute.String("http.path", string(ctx.Path())),
			attribute.String("http.url", string(ctx.URI().FullURI())),
			attribute.String("http.user_agent", string(ctx.UserAgent())),
			attribute.String("http.remote_addr", ctx.RemoteAddr().String()),
		)

		// Add span to fasthttp context
		ctx.SetUserValue("otel_span", span)
		ctx.SetUserValue("otel_context", reqCtx)

		// Track execution time
		start := time.Now()

		// Execute main handler
		handler(ctx)

		// Record duration and response status
		duration := time.Since(start)
		statusCode := ctx.Response.StatusCode()

		// Add response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.String("http.status_text", fasthttp.StatusMessage(statusCode)),
			attribute.Int64("http.response_size", int64(ctx.Response.Header.ContentLength())),
		)

		// Record operation duration
		telemetry.RecordSpanDuration(span, fmt.Sprintf("HTTP %s", string(ctx.Method())), duration)

		// Add request completion event
		telemetry.AddSpanEvent(span, "request_completed",
			attribute.Int("status_code", statusCode),
			attribute.Int64("duration_ms", duration.Milliseconds()),
		)

		// If there's an error (status >= 400), record it
		if statusCode >= 400 {
			errorMsg := fmt.Sprintf("HTTP %d: %s", statusCode, fasthttp.StatusMessage(statusCode))
			telemetry.SetSpanError(span, fmt.Errorf(errorMsg))
		}

		// End span
		telemetry.EndSpan(span)
	}
}

// GetSpanFromContext extracts span from fasthttp context
func GetSpanFromContext(ctx *fasthttp.RequestCtx) trace.Span {
	if spanValue := ctx.UserValue("otel_span"); spanValue != nil {
		if span, ok := spanValue.(trace.Span); ok {
			return span
		}
	}
	return trace.SpanFromContext(context.Background())
}

// GetOtelContextFromRequest extracts OpenTelemetry context from request
func GetOtelContextFromRequest(ctx *fasthttp.RequestCtx) context.Context {
	if ctxValue := ctx.UserValue("otel_context"); ctxValue != nil {
		if otelCtx, ok := ctxValue.(context.Context); ok {
			return otelCtx
		}
	}
	return context.Background()
}

// AddSpanAttributes adds attributes to current span
func AddSpanAttributes(ctx *fasthttp.RequestCtx, attrs ...attribute.KeyValue) {
	span := GetSpanFromContext(ctx)
	if span != nil {
		span.SetAttributes(attrs...)
	}
}

// AddSpanEventToRequest adds event to current span
func AddSpanEventToRequest(ctx *fasthttp.RequestCtx, name string, attrs ...attribute.KeyValue) {
	span := GetSpanFromContext(ctx)
	if span != nil {
		telemetry.AddSpanEvent(span, name, attrs...)
	}
}

// RecordSpanError records error in current span
func RecordSpanError(ctx *fasthttp.RequestCtx, err error) {
	span := GetSpanFromContext(ctx)
	if span != nil {
		telemetry.SetSpanError(span, err)
	}
}

// CreateChildSpan creates child span for additional operations
func CreateChildSpan(ctx *fasthttp.RequestCtx, operationName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	parentCtx := GetOtelContextFromRequest(ctx)
	tracer := telemetry.GetTracer("k8s-controller-tutorial/api")

	return telemetry.StartSpan(parentCtx, tracer, operationName, attrs...)
}

// TraceableHandler wraps handler with additional tracing parameters
func TraceableHandler(operationName string, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Add additional attributes for specific operation
		AddSpanAttributes(ctx,
			attribute.String("operation.name", operationName),
			attribute.String("handler.type", "api"),
		)

		// Add operation start event
		AddSpanEventToRequest(ctx, fmt.Sprintf("start_%s", strings.ToLower(operationName)))

		// Execute main handler
		handler(ctx)

		// Add operation end event
		AddSpanEventToRequest(ctx, fmt.Sprintf("end_%s", strings.ToLower(operationName)))
	}
}
