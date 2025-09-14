package memory

import (
	"sync"
	"time"

	"github.com/TheVidz/Kafkaesque/pkg/broker"
)

type MemoryBroker struct {
	topics map[string]*broker.Topic
	mu     sync.RWMutex
}

func NewMemoryBroker() *MemoryBroker {
	return &MemoryBroker{
		topics: make(map[string]*broker.Topic),
	}
}

func (mb *MemoryBroker) getOrCreateTopic(topic string) *broker.Topic {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	if t, ok := mb.topics[topic]; ok {
		return t
	}
	newTopic := broker.NewTopic(topic)
	mb.topics[topic] = newTopic
	return newTopic
}

func (mb *MemoryBroker) Publish(topic string, payload []byte) error {
	t := mb.getOrCreateTopic(topic)

	msg := broker.Message{
		ID:        time.Now().UnixNano(),
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	// Phase-2 publish
	t.Publish(msg)
	return nil
}

// Instead of taking a raw channel, we now take a *Subscriber
func (mb *MemoryBroker) Subscribe(topic string, sub *broker.Subscriber) error {
	t := mb.getOrCreateTopic(topic)
	t.AddSubscriber(sub)
	return nil
}
