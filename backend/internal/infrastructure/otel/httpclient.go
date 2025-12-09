package otel

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Transport that captures request and response bodies
type bodyCapturingTransport struct {
	http.RoundTripper
}

func (t *bodyCapturingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	span := trace.SpanFromContext(req.Context())

	// Always execute the request, but only capture bodies if enabled
	if !span.IsRecording() {
		return t.RoundTripper.RoundTrip(req)
	}

	// Capture REQUEST BODY only if enabled
	if shouldCaptureBody() && req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > 0 && isTextContent(req.Header.Get("Content-Type")) {
				bodyString := string(bodyBytes)
				bodyString = strings.TrimSpace(bodyString)
				span.SetAttributes(attribute.String("http.request.body", bodyString))
			}
		}
	}

	// Make the request
	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		// Marcar como error si hay error de red/conexión
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.Bool("error", true))
		return resp, err
	}

	// Capture RESPONSE BODY only if enabled
	if shouldCaptureBody() && resp.Body != nil {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if readErr == nil {
			// Restaurar el body para el consumidor original
			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > 0 && isTextContent(resp.Header.Get("Content-Type")) {
				bodyString := string(bodyBytes)
				bodyString = strings.TrimSpace(bodyString)
				span.SetAttributes(attribute.String("http.response.body", bodyString))
			}
		}
	}

	// MARK AS ERROR IF STATUS CODE >= 400
	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", resp.StatusCode))
		span.SetAttributes(
			attribute.Bool("error", true),
			attribute.String("error.type", "http_error"),
			attribute.Int("http.status_code", resp.StatusCode),
		)

		// Add more descriptive error message based on the range
		var errorMessage string
		switch {
		case resp.StatusCode >= 500:
			errorMessage = "Server Error"
		case resp.StatusCode >= 400:
			errorMessage = "Client Error"
		}
		span.SetAttributes(attribute.String("error.message", errorMessage))
	} else {
		// Mark as success explicitly
		span.SetStatus(codes.Ok, "")
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	}

	return resp, err
}

func SetupHttpClient() {
	// Create transport with body capture
	baseTransport := &bodyCapturingTransport{
		RoundTripper: http.DefaultTransport,
	}

	// Wrap with otelhttp for base instrumentation
	transport := otelhttp.NewTransport(
		baseTransport,
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String())
		}),
	)

	http.DefaultTransport = transport
	http.DefaultClient.Transport = transport
}

// Check if bodies should be captured based on environment variable
func shouldCaptureBody() bool {
	return os.Getenv("OTEL_CAPTURE_BODIES") == "true"
}

// Check if content is text to avoid binary data
func isTextContent(contentType string) bool {
	textTypes := []string{
		"application/json",
		"application/xml",
		"text/",
		"application/x-www-form-urlencoded",
	}

	contentType = strings.ToLower(contentType)
	for _, textType := range textTypes {
		if strings.Contains(contentType, textType) {
			return true
		}
	}
	return false
}
