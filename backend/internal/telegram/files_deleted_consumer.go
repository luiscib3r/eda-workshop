package telegram

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/storage/events"
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type FilesDeletedConsumer struct {
	*nats.NatsConsumer[*events.FilesDeletedEvent]
	bot *TelegramBot
}

func NewFilesDeletedConsumer(
	js jetstream.JetStream,
	bot *TelegramBot,
) *FilesDeletedConsumer {
	name := "tgbot_files_deleted_consumer"

	numWorkers := 4
	workerBufferSize := 10
	consumer := &FilesDeletedConsumer{
		bot: bot,
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
			Description:   "Telegram Bot Files Deleted Event Consumer",
			FilterSubject: events.STORAGE_FILES_DELETED_EVENT,
		},
	)

	return consumer
}

func (c *FilesDeletedConsumer) handler(
	ctx context.Context,
	event *events.FilesDeletedEvent,
) error {

	msg := "🗑️ *Files Deleted!*\n\n"
	msg += "The following files have been deleted:\n\n"

	for _, fileKey := range event.Payload.FileKeys {
		msg += "• `" + fileKey + "`\n"
	}

	if err := c.bot.SendMessage(msg); err != nil {
		return err
	}

	return nil
}
