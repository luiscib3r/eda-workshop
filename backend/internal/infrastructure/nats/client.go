package nats

import (
	"backend/internal/infrastructure/config"
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

func NewNatsClient(
	lc fx.Lifecycle,
	cfg *config.AppConfig,
) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.Nats.Uri)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if nc != nil && !nc.IsClosed() {
				nc.Close()
			}
			return nil
		},
	})

	return nc, nil
}

func NewJetStreamClient(
	nc *nats.Conn,
) (jetstream.JetStream, error) {
	return jetstream.New(nc)
}
