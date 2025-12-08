package events

import (
	"backend/gen/ocr"
	"backend/internal/core"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
)

type FilePageRegisteredEvent struct {
	Id      ulid.ULID
	Payload *ocr.FilePageRegisteredEventData
}

var _ core.EventSpec = (*FilePageRegisteredEvent)(nil)

func NewFilePageRegisteredEvent(
	payload *ocr.FilePageRegisteredEventData,
) *FilePageRegisteredEvent {
	return &FilePageRegisteredEvent{
		Id:      core.NewEventID(),
		Payload: payload,
	}
}

func NewFilePageRegisteredEventFromMessage(
	msg jetstream.Msg,
) (*FilePageRegisteredEvent, error) {
	headers := msg.Headers()
	data := msg.Data()

	payload := &ocr.FilePageRegisteredEventData{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	id, err := ulid.Parse(headers.Get(core.EVENT_ID_HEADER))
	if err != nil {
		return nil, err
	}

	event := &FilePageRegisteredEvent{
		Id:      id,
		Payload: payload,
	}

	return event, nil
}

// ID implements core.EventSpec.
func (ev *FilePageRegisteredEvent) ID() string {
	return ev.Id.String()
}

// Type implements core.EventSpec.
func (ev *FilePageRegisteredEvent) Type() string {
	return FILE_PAGE_REGISTERED_EVENT
}

// Data implements core.EventSpec.
func (ev *FilePageRegisteredEvent) Data() proto.Message {
	return ev.Payload
}
