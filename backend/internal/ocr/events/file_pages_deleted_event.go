package events

import (
	"backend/gen/ocr"
	"backend/internal/core"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
)

type FilePagesDeletedEvent struct {
	Id      ulid.ULID
	Payload *ocr.FilePagesDeletedEventData
}

var _ core.EventSpec = (*FilePagesDeletedEvent)(nil)

func NewFilePagesDeletedEvent(
	payload *ocr.FilePagesDeletedEventData,
) *FilePagesDeletedEvent {
	return &FilePagesDeletedEvent{
		Id:      core.NewEventID(),
		Payload: payload,
	}
}

func NewFilePagesDeletedEventFromMessage(
	msg jetstream.Msg,
) (*FilePagesDeletedEvent, error) {
	headers := msg.Headers()
	data := msg.Data()

	payload := &ocr.FilePagesDeletedEventData{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	id, err := ulid.Parse(headers.Get(core.EVENT_ID_HEADER))
	if err != nil {
		return nil, err
	}

	event := &FilePagesDeletedEvent{
		Id:      id,
		Payload: payload,
	}

	return event, nil
}

// ID implements core.EventSpec.
func (ev *FilePagesDeletedEvent) ID() string {
	return ev.Id.String()
}

// Type implements core.EventSpec.
func (ev *FilePagesDeletedEvent) Type() string {
	return FILE_PAGES_DELETED_EVENT
}

// Data implements core.EventSpec.
func (ev *FilePagesDeletedEvent) Data() proto.Message {
	return ev.Payload
}
