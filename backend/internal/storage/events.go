package storage

import (
	"backend/gen/storage"
	"backend/internal/core"

	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

const (
	STORAGE_CHANNEL             string = "storage"
	STORAGE_FILE_UPLOADED_EVENT string = "storage.file.uploaded"
)

type FileUploadedEvent struct {
	id      string
	Payload *storage.FileUploadedEventData
}

var _ core.EventSpec = (*FileUploadedEvent)(nil)

func NewFileUploadedEvent(
	payload *storage.FileUploadedEventData,
) *FileUploadedEvent {
	return &FileUploadedEvent{
		id:      core.NewEventID(),
		Payload: payload,
	}
}

func NewFileUploadedEventFromMessage(
	msg jetstream.Msg,
) (*FileUploadedEvent, error) {
	headers := msg.Headers()
	data := msg.Data()

	payload := &storage.FileUploadedEventData{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	event := &FileUploadedEvent{
		id:      headers.Get("Event-ID"),
		Payload: payload,
	}

	return event, nil
}

// ID implements core.EventSpec.
func (ev *FileUploadedEvent) ID() string {
	return ev.id
}

// Type implements core.EventSpec.
func (ev *FileUploadedEvent) Type() string {
	return "storage.file.uploaded"
}

// Data implements core.EventSpec.
func (ev *FileUploadedEvent) Data() proto.Message {
	return ev.Payload
}
