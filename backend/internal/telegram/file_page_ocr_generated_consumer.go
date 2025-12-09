package telegram

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/ocr/events"
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type FilePageOcrGenerateConsumer struct {
	*nats.NatsConsumer[*events.FilePageOcrGeneratedEvent]
	bot *TelegramBot
}

func NewFilePageOcrGenerateConsumer(
	js jetstream.JetStream,
	bot *TelegramBot,
) *FilePageOcrGenerateConsumer {
	name := "tgbot_file_page_ocr_generated_consumer"
	numWorkers := 4
	workerBufferSize := 10
	consumer := &FilePageOcrGenerateConsumer{
		bot: bot,
	}

	consumer.NatsConsumer = nats.NewNatsConsumer(
		name,
		events.OCR_CHANNEL,
		events.FILE_PAGE_OCR_GENERATED_EVENT,
		numWorkers,
		workerBufferSize,
		events.NewFilePageOcrGeneratedEventFromMessage,
		consumer.handler,
		js,
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			Description:   "Telegram Bot File Page OCR Generated Event Consumer",
			FilterSubject: events.FILE_PAGE_OCR_GENERATED_EVENT,
		},
	)

	return consumer
}

func (c *FilePageOcrGenerateConsumer) handler(
	ctx context.Context,
	event *events.FilePageOcrGeneratedEvent,
) error {
	msg := "📝 *New OCR Generated!*\n\n"
	msg += "😎 *File Key: " + event.Payload.Id + "*\n"
	msg += "🔑 *Page Image Key:* `" + event.Payload.PageImageKey + "`\n"
	msg += "📄 *Page Number:* `" + fmt.Sprintf("%d", event.Payload.PageNumber) + "`\n"
	msg += "\n✅ OCR generated successfully!"

	if err := c.bot.SendMessage(msg); err != nil {
		return err
	}

	return nil
}
