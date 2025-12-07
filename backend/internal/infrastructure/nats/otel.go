package nats

import (
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/propagation"
)

type natsCarrier struct {
	headers nats.Header
}

var _ propagation.TextMapCarrier = (*natsCarrier)(nil)

// Get implements propagation.TextMapCarrier.
func (n *natsCarrier) Get(key string) string {
	return n.headers.Get(key)
}

// Set implements propagation.TextMapCarrier.
func (n *natsCarrier) Set(key string, value string) {
	n.headers.Set(key, value)
}

// Keys implements propagation.TextMapCarrier.
func (n *natsCarrier) Keys() []string {
	keys := make([]string, 0, len(n.headers))
	for key := range n.headers {
		keys = append(keys, key)
	}
	return keys
}
