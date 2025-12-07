package events

import (
	"backend/gen/storage"
	"backend/internal/core"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
)

type FileUploadedEvent struct {
	Id      ulid.ULID
	Payload *storage.FileUploadedEventData
}

var _ core.EventSpec = (*FileUploadedEvent)(nil)

func NewFileUploadedEvent(
	payload *storage.FileUploadedEventData,
) *FileUploadedEvent {
	return &FileUploadedEvent{
		Id:      core.NewEventID(),
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

	id, err := ulid.Parse(headers.Get(core.EVENT_ID_HEADER))
	if err != nil {
		return nil, err
	}

	event := &FileUploadedEvent{
		Id:      id,
		Payload: payload,
	}

	return event, nil
}

// ID implements core.EventSpec.
func (ev *FileUploadedEvent) ID() string {
	return ev.Id.String()
}

// Type implements core.EventSpec.
func (ev *FileUploadedEvent) Type() string {
	return STORAGE_FILE_UPLOADED_EVENT
}

// Data implements core.EventSpec.
func (ev *FileUploadedEvent) Data() proto.Message {
	return ev.Payload
}
