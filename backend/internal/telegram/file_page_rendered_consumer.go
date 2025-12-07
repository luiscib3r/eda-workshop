package telegram

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/ocr/events"
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type FilePageRenderedConsumer struct {
	*nats.NatsConsumer[*events.FilePageRenderedEvent]
	bot *TelegramBot
}

func NewFilePageRenderedConsumer(
	js jetstream.JetStream,
	bot *TelegramBot,
) *FilePageRenderedConsumer {
	name := "tgbot_file_page_rendered_consumer"
	numWorkers := 5
	workerBufferSize := 10
	consumer := &FilePageRenderedConsumer{
		bot: bot,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.OCR_CHANNEL,
		events.FILE_PAGE_RENDERED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFilePageRenderedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "Telegram Bot File Page Rendered Event Consumer",
			FilterSubject: events.FILE_PAGE_RENDERED_EVENT,
		},
	)

	return consumer
}

func (c *FilePageRenderedConsumer) handler(
	ctx context.Context,
	event *events.FilePageRenderedEvent,
) error {
	msg := "🖼 *New Page Rendered!*\n\n"
	msg += "😎 *File Key: " + event.Payload.FileKey + "*\n"
	msg += "📄 *Page Image Key:* `" + event.Payload.PageImageKey + "`\n"
	msg += "\n✅ Page rendered successfully!"

	if err := c.bot.SendMessage(msg); err != nil {
		return err
	}

	return nil
}
