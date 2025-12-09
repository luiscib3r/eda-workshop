package events

import (
	"backend/gen/ocr"
	"backend/internal/core"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
)

type FilePageOcrGeneratedEvent struct {
	Id      ulid.ULID
	Payload *ocr.FilePageOcrGeneratedEventData
}

var _ core.EventSpec = (*FilePageOcrGeneratedEvent)(nil)

func NewFilePageOcrGeneratedEvent(
	payload *ocr.FilePageOcrGeneratedEventData,
) *FilePageOcrGeneratedEvent {
	return &FilePageOcrGeneratedEvent{
		Id:      core.NewEventID(),
		Payload: payload,
	}
}

func NewFilePageOcrGeneratedEventFromMessage(
	msg jetstream.Msg,
) (*FilePageOcrGeneratedEvent, error) {
	headers := msg.Headers()
	data := msg.Data()

	payload := &ocr.FilePageOcrGeneratedEventData{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	id, err := ulid.Parse(headers.Get(core.EVENT_ID_HEADER))
	if err != nil {
		return nil, err
	}

	event := &FilePageOcrGeneratedEvent{
		Id:      id,
		Payload: payload,
	}

	return event, nil
}

// ID implements core.EventSpec.
func (ev *FilePageOcrGeneratedEvent) ID() string {
	return ev.Id.String()
}

// Type implements core.EventSpec.
func (ev *FilePageOcrGeneratedEvent) Type() string {
	return FILE_PAGE_OCR_GENERATED_EVENT
}

// Data implements core.EventSpec.
func (ev *FilePageOcrGeneratedEvent) Data() proto.Message {
	return ev.Payload
}
