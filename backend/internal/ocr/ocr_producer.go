package ocr

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/ocr/events"

	"github.com/nats-io/nats.go/jetstream"
)

type OcrProducer struct {
	*nats.NatsProducer
}

func NewOcrProducer(
	js jetstream.JetStream,
) *OcrProducer {
	name := "ocr_producer"

	return &OcrProducer{
		NatsProducer: nats.NewNatsProducer(
			// Producer Name
			name,
			// Channel Name
			events.OCR_CHANNEL,
			// JetStream Context
			js,
			// Stream Configuration
			jetstream.StreamConfig{
				Name:        nats.StreamName(events.OCR_CHANNEL),
				Description: "OCR Service Event Stream",
				Subjects:    []string{events.OCR_CHANNEL + ".>"},
			},
		),
	}
}
