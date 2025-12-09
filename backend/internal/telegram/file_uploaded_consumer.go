package telegram

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/storage/events"
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type FileUploadedConsumer struct {
	*nats.NatsConsumer[*events.FileUploadedEvent]
	bot *TelegramBot
}

func NewFileUploadedConsumer(
	js jetstream.JetStream,
	bot *TelegramBot,
) *FileUploadedConsumer {
	name := "tgbot_file_uploaded_consumer"

	numWorkers := 4
	workerBufferSize := 10
	consumer := &FileUploadedConsumer{
		bot: bot,
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
			Description:   "Telegram Bot File Uploaded Event Consumer",
			FilterSubject: events.STORAGE_FILE_UPLOADED_EVENT,
		},
	)

	return consumer
}

func (c *FileUploadedConsumer) handler(
	ctx context.Context,
	event *events.FileUploadedEvent,
) error {

	msg := "📤 *New File Uploaded!*\n\n"
	msg += "😎 *File Name: " + event.Payload.FileName + "*\n"
	msg += "📄 *File Key:* `" + event.Payload.FileKey + "`\n"
	msg += "\n✅ Upload completed successfully!"

	if err := c.bot.SendMessage(msg); err != nil {
		return err
	}

	return nil
}
