package storage

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/storage/events"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/nats-io/nats.go/jetstream"
)

type FilesDeletedConsumer struct {
	*nats.NatsConsumer[*events.FilesDeletedEvent]
	s3 *s3.Client
}

func NewFilesDeletedConsumer(
	js jetstream.JetStream,
	s3 *s3.Client,
) *FilesDeletedConsumer {
	name := "storage_files_deleted_consumer"
	numWorkers := 5
	workerBufferSize := 10

	consumer := &FilesDeletedConsumer{
		s3: s3,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.STORAGE_CHANNEL,
		events.STORAGE_FILES_DELETED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFilesDeletedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "Storage Files Deleted Event Consumer",
			FilterSubject: events.STORAGE_FILES_DELETED_EVENT,
		},
	)

	return consumer
}

func (c *FilesDeletedConsumer) handler(
	ctx context.Context,
	event *events.FilesDeletedEvent,
) error {
	bucket := "files"
	fileKeys := event.Payload.FileKeys

	// Build the list of objects to delete
	objects := make([]types.ObjectIdentifier, len(fileKeys))
	for i, fileKey := range fileKeys {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(fileKey),
		}
	}

	// Delete multiple files from S3 in a single call
	_, err := c.s3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})

	return err
}
