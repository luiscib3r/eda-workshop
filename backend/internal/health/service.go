package health

import (
	"backend/gen/health"
	"backend/internal/infrastructure/service"
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HealthService struct {
	health.UnimplementedHealthServiceServer
}

var _ health.HealthServiceServer = (*HealthService)(nil)
var _ service.Service = (*HealthService)(nil)

func NewHealthService() *HealthService {
	return &HealthService{}
}

// Health implements health.HealthServiceServer.
func (h *HealthService) Health(context.Context, *emptypb.Empty) (*health.HealthResponse, error) {
	return &health.HealthResponse{
		Status: "OK",
	}, nil
}

// Register implements service.Service.
func (h *HealthService) Register(ctx context.Context, mux *runtime.ServeMux) {
	health.RegisterHealthServiceHandlerServer(ctx, mux, h)
}
