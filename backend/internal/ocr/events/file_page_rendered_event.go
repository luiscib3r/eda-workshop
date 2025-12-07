package events

import (
	"backend/gen/ocr"
	"backend/internal/core"

	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

type FilePageRenderedEvent struct {
	Id      string
	Payload *ocr.FilePageRenderedEventData
}

var _ core.EventSpec = (*FilePageRenderedEvent)(nil)

func NewFilePageRenderedEvent(
	payload *ocr.FilePageRenderedEventData,
) *FilePageRenderedEvent {
	return &FilePageRenderedEvent{
		Id:      core.NewEventID(),
		Payload: payload,
	}
}

func NewFilePageRenderedEventFromMessage(
	msg jetstream.Msg,
) (*FilePageRenderedEvent, error) {
	headers := msg.Headers()
	data := msg.Data()

	payload := &ocr.FilePageRenderedEventData{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, err
	}

	event := &FilePageRenderedEvent{
		Id:      headers.Get("Event-ID"),
		Payload: payload,
	}

	return event, nil
}

// ID implements core.EventSpec.
func (ev *FilePageRenderedEvent) ID() string {
	return ev.Id
}

// Type implements core.EventSpec.
func (ev *FilePageRenderedEvent) Type() string {
	return FILE_PAGE_RENDERED_EVENT
}

// Data implements core.EventSpec.
func (ev *FilePageRenderedEvent) Data() proto.Message {
	return ev.Payload
}
