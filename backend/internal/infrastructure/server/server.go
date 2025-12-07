package server

import (
	"backend/internal/infrastructure/config"
	"backend/openapi"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewGatewayServeMux() *runtime.ServeMux {
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(OtelTraceIDHeader),
	)

	return mux
}

func NewServeMux(
	gateway *runtime.ServeMux,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Setup docs
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.FS(openapi.OpenApiFS))))

	// Setup gRPC-Gateway
	mux.Handle("/", gateway)

	return mux
}

func NewHttpServer(
	cfg *config.AppConfig,
	mux *http.ServeMux,
) *http.Server {
	handler := otelhttp.NewHandler(
		PanicRecover(mux),
		"grpc-gateway",
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		}),
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port),
		Handler: cors(cfg.Cors)(handler),
	}

	return srv
}
