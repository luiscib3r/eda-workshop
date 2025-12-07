package storage

import (
	"backend/internal/infrastructure/nats"
	"backend/internal/storage/events"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

type StorageProducer struct {
	*nats.NatsProducer
}

func NewStorageProducer(
	js jetstream.JetStream,
) *StorageProducer {
	name := "storage_producer"

	return &StorageProducer{
		NatsProducer: nats.NewNatsProducer(
			// Producer Name
			name,
			// Channel Name
			events.STORAGE_CHANNEL,
			// JetStream Context
			js,
			// Stream Configuration
			jetstream.StreamConfig{
				Name:        nats.StreamName(events.STORAGE_CHANNEL),
				Description: "Storage Service Event Stream",
				Subjects:    []string{fmt.Sprintf("%s.>", events.STORAGE_CHANNEL)},
			},
		),
	}
}
