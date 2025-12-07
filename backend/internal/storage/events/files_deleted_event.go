package events

import (
	"backend/gen/storage"
	"backend/internal/core"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
)

type FilesDeletedEvent struct {
	Id      ulid.ULID
	Payload *storage.FilesDeletedEventData
}

var _ core.EventSpec = (*FilesDeletedEvent)(nil)

func NewFilesDeletedEvent(
	payload *storage.FilesDeletedEventData,
) *FilesDeletedEvent {
	return &FilesDeletedEvent{
		Id:      core.NewEventID(),
		Payload: payload,
	}
}

func NewFilesDeletedEventFromMessage(
	msg jetstream.Msg,
) (*FilesDeletedEvent, error) {
	headers := msg.Headers()
	data := msg.Data()

	payload := &storage.FilesDeletedEventData{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	id, err := ulid.Parse(headers.Get(core.EVENT_ID_HEADER))
	if err != nil {
		return nil, err
	}

	event := &FilesDeletedEvent{
		Id:      id,
		Payload: payload,
	}

	return event, nil
}

// ID implements core.EventSpec.
func (ev *FilesDeletedEvent) ID() string {
	return ev.Id.String()
}

// Type implements core.EventSpec.
func (ev *FilesDeletedEvent) Type() string {
	return STORAGE_FILES_DELETED_EVENT
}

// Data implements core.EventSpec.
func (ev *FilesDeletedEvent) Data() proto.Message {
	return ev.Payload
}
