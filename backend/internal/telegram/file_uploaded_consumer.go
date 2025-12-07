package telegram

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/storage"
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type FileUploadedConsumer struct {
	*nats.NatsConsumer[*storage.FileUploadedEvent]
	bot *TelegramBot
}

func NewFileUploadedConsumer(
	js jetstream.JetStream,
	bot *TelegramBot,
) *FileUploadedConsumer {
	name := "tgbot_file_uploaded_consumer"

	numWorkers := 5
	workerBufferSize := 10
	consumer := &FileUploadedConsumer{
		bot: bot,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		storage.STORAGE_CHANNEL,
		storage.STORAGE_FILE_UPLOADED_EVENT,
		numWorkers,
		workerBufferSize,
		storage.NewFileUploadedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "Telegram Bot File Uploaded Event Consumer",
			FilterSubject: storage.STORAGE_FILE_UPLOADED_EVENT,
		},
	)

	return consumer
}

func (c *FileUploadedConsumer) handler(
	ctx context.Context,
	event *storage.FileUploadedEvent,
) error {

	msg := "📤 *New File Uploaded!*\n\n"
	msg += "😎 *File Name: " + event.Payload.FileName + "*\n"
	msg += "📄 *File Key:* `" + event.Payload.FileKey + "`\n"
	msg += "🗂️ *Bucket Name:* `" + event.Payload.BucketName + "`\n"
	msg += "\n✅ Upload completed successfully!"

	if err := c.bot.SendMessage(msg); err != nil {
		return err
	}

	return nil
}
