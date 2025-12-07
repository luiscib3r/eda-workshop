package telegram

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/storage"
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type FileUploadedEventConsumer struct {
	*nats.NatsConsumer[*storage.FileUploadedEvent]
	bot *TelegramBot
}

func NewFileUploadedEventConsumer(
	js jetstream.JetStream,
	bot *TelegramBot,
) *FileUploadedEventConsumer {
	name := "tgbot_file_uploaded_consumer"
	channel := "storage"
	event := "storage.file.uploaded"
	numWorkers := 5
	workerBufferSize := 10
	consumer := &FileUploadedEventConsumer{
		bot: bot,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		channel,
		event,
		numWorkers,
		workerBufferSize,
		storage.NewFileUploadedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "Telegram Bot File Uploaded Event Consumer",
			FilterSubject: event,
		},
	)

	return consumer
}

func (c *FileUploadedEventConsumer) handler(
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
