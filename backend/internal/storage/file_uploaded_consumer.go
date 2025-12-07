package storage

import (
	"backend/internal/infrastructure/nats"
	storagedb "backend/internal/storage/db"
	"backend/internal/storage/events"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nats-io/nats.go/jetstream"
)

type FileUploadedConsumer struct {
	*nats.NatsConsumer[*events.FileUploadedEvent]
	db *storagedb.Queries
	s3 *s3.Client
}

func NewFileUploadedConsumer(
	js jetstream.JetStream,
	db *storagedb.Queries,
	s3 *s3.Client,
) *FileUploadedConsumer {
	name := "storage_file_uploaded_consumer"
	numWorkers := 10
	workerBufferSize := 10

	consumer := &FileUploadedConsumer{
		db: db,
		s3: s3,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.STORAGE_CHANNEL,
		events.STORAGE_FILE_UPLOADED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFileUploadedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "Storage File Uploaded Event Consumer",
			FilterSubject: events.STORAGE_FILE_UPLOADED_EVENT,
		},
	)

	return consumer
}

func (c *FileUploadedConsumer) handler(
	ctx context.Context,
	event *events.FileUploadedEvent,
) error {
	// Get file info
	head, err := c.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(event.Payload.BucketName),
		Key:    aws.String(event.Payload.FileKey),
	})

	if err != nil {
		return err
	}

	var size int64 = 0
	if head.ContentLength != nil {
		size = int64(*head.ContentLength)
	}
	fileType := ""
	if head.ContentType != nil {
		fileType = *head.ContentType
	}

	// Create file record
	if _, err := c.db.CreateFile(ctx, storagedb.CreateFileParams{
		ID:         event.Payload.FileKey,
		FileName:   event.Payload.FileName,
		BucketName: event.Payload.BucketName,
		FileSize:   size,
		FileType:   fileType,
	}); err != nil {
		return err
	}

	return nil
}
