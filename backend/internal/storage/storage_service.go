package storage

import (
	"backend/gen/storage"
	"backend/internal/infrastructure/config"
	"backend/internal/infrastructure/service"
	"backend/internal/storage/events"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageService struct {
	storage.UnimplementedStorageServiceServer
	presign  *s3.PresignClient
	client   *s3.Client
	producer *StorageProducer
	cors     config.CorsConfig
}

var _ storage.StorageServiceServer = (*StorageService)(nil)
var _ service.Service = (*StorageService)(nil)

func NewStorageService(
	presign *s3.PresignClient,
	client *s3.Client,
	producer *StorageProducer,
	cfg *config.AppConfig,
) *StorageService {
	return &StorageService{
		presign:  presign,
		client:   client,
		producer: producer,
		cors:     cfg.Cors,
	}
}

// GetUploadUrl implements storage.StorageServiceServer.
func (s *StorageService) GetUploadUrl(
	ctx context.Context,
	req *emptypb.Empty,
) (*storage.GetUploadUrlResponse, error) {
	// Generate random object key
	bucket := "files"
	key := ulid.MustNew(
		ulid.Timestamp(time.Now()),
		ulid.DefaultEntropy(),
	).String()

	// Generate presigned URL
	result, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Key:    &key,
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}

	return &storage.GetUploadUrlResponse{
		UploadUrl: result.URL,
		FileKey:   key,
	}, nil
}

// ConfirmFileUpload implements storage.StorageServiceServer.
func (s *StorageService) ConfirmFileUpload(
	ctx context.Context,
	req *storage.ConfirmFileUploadRequest,
) (*emptypb.Empty, error) {
	event := events.NewFileUploadedEvent(
		&storage.FileUploadedEventData{
			FileName: req.FileName,
			FileKey:  req.FileKey,
		},
	)

	err := s.producer.Publish(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to publish file uploaded event: %w", err)
	}

	return &emptypb.Empty{}, nil
}

// Register implements service.Service.
func (s *StorageService) Register(ctx context.Context, mux *runtime.ServeMux) {
	storage.RegisterStorageServiceHandlerServer(ctx, mux, s)
}

func (s *StorageService) CreateBuckets(ctx context.Context) error {
	buckets := []string{"files"}

	for _, bucket := range buckets {
		_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucket),
		})

		if err == nil {
			slog.InfoContext(ctx, "Bucket already exists", "bucket", bucket)
			continue
		}

		_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create bucket", "bucket", bucket, "error", err)
			var bae *types.BucketAlreadyExists
			var baoy *types.BucketAlreadyOwnedByYou
			if !errors.As(err, &bae) && !errors.As(err, &baoy) {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}

		_, err = s.client.PutBucketCors(ctx, &s3.PutBucketCorsInput{
			Bucket: aws.String(bucket),
			CORSConfiguration: &types.CORSConfiguration{
				CORSRules: []types.CORSRule{
					{
						AllowedOrigins: s.cors.AllowedOrigins,
						AllowedHeaders: []string{"*"},
						AllowedMethods: []string{"*"},
						MaxAgeSeconds:  aws.Int32(3000),
					},
				},
			},
		})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to set CORS for bucket", "bucket", bucket, "error", err)
		}
	}

	return nil
}
