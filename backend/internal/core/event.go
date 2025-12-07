package core

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
)

type EventSpec interface {
	ID() string
	Type() string
	Data() proto.Message
}

func NewEventID() string {
	return ulid.MustNew(
		ulid.Timestamp(time.Now()),
		ulid.DefaultEntropy(),
	).String()
}

type EventBuilder[T EventSpec] func(msg jetstream.Msg) (T, error)
type EventHandler[T EventSpec] func(ctx context.Context, event T) error
