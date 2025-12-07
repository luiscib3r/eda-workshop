package storage

import (
	"backend/internal/infrastructure/config"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

const BUCKET_NAME = "files"

func NewClient(
	cfg *config.AppConfig,
) (*s3.Client, error) {
	cred := credentials.NewStaticCredentialsProvider(
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		"",
	)

	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(cfg.Storage.Region),
		awsConfig.WithCredentialsProvider(cred),
	)
	if err != nil {
		return nil, err
	}
	otelaws.AppendMiddlewares(&awsCfg.APIOptions)

	client := s3.NewFromConfig(awsCfg, func(opts *s3.Options) {
		opts.BaseEndpoint = aws.String(cfg.Storage.Endpoint)
		opts.UsePathStyle = cfg.Storage.UsePathStyle
	})

	return client, nil
}

func NewPresignClient(
	cfg *config.AppConfig,
) (*s3.PresignClient, error) {
	cred := credentials.NewStaticCredentialsProvider(
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		"",
	)

	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(cfg.Storage.Region),
		awsConfig.WithCredentialsProvider(cred),
	)
	if err != nil {
		return nil, err
	}
	otelaws.AppendMiddlewares(&awsCfg.APIOptions)

	client := s3.NewFromConfig(awsCfg, func(opts *s3.Options) {
		opts.BaseEndpoint = aws.String(cfg.Storage.PublicEndpoint)
		opts.UsePathStyle = cfg.Storage.UsePathStyle
	})

	return s3.NewPresignClient(client), nil
}
