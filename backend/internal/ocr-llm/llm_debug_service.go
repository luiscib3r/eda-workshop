package ocrllm

import (
	"backend/gen/ocr"
	"backend/internal/infrastructure/service"
	"backend/internal/infrastructure/storage"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type LlmDebugService struct {
	ocr.UnimplementedLlmDebugServiceServer
	ocr *OcrAgent
	s3  *s3.Client
}

var _ ocr.LlmDebugServiceServer = (*LlmDebugService)(nil)
var _ service.Service = (*LlmDebugService)(nil)

func NewLlmDebugService(
	ocr *OcrAgent,
	s3 *s3.Client,
) *LlmDebugService {
	return &LlmDebugService{
		ocr: ocr,
		s3:  s3,
	}
}

// GetOcr implements ocr.LlmDebugServiceServer.
func (l *LlmDebugService) GetOcr(
	ctx context.Context,
	req *ocr.GetOcrRequest,
) (*ocr.GetOcrResponse, error) {
	// Fetch image from S3
	result, err := l.s3.GetObject(ctx, &s3.GetObjectInput{
		Key:    &req.PageKey,
		Bucket: aws.String(storage.BUCKET_NAME),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	resp, err := l.ocr.Invoke(ctx, data)
	if err != nil {
		return nil, err
	}

	ocrText := "Not recognized"
	if len(resp.Choices) > 0 {
		ocrText = resp.Choices[0].Message.Content
	}

	return &ocr.GetOcrResponse{
		OcrText: ocrText,
	}, nil
}

// Register implements service.Service.
func (l *LlmDebugService) Register(ctx context.Context, mux *runtime.ServeMux) {
	ocr.RegisterLlmDebugServiceHandlerServer(ctx, mux, l)
}
