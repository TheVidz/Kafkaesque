package broker

import (
	"context"
	"time"
)

type DeliveryGuarantee int

const (
	AtMostOnce DeliveryGuarantee = iota
	AtLeastOnce
)

type TopicOptions struct {
	Partitions int
	Retention  time.Duration
	MaxBytes   int64
}

type Message struct {
	Topic     string
	Key       string
	Value     []byte
	Offset    int64
	Timestamp time.Time
	Headers   map[string]string
}

type Producer interface {
	Publish(ctx context.Context, topic, key string, value []byte) (int64, error)
}

type Consumer interface {
	// Subscribe returns a channel of messages and a stop func to cancel the subscription.
	Subscribe(ctx context.Context, topic, group string) (<-chan Message, func(), error)
}

type Broker interface {
	CreateTopic(ctx context.Context, name string, opts TopicOptions) error
	Producer() Producer
	Consumer() Consumer
	Close() error
}