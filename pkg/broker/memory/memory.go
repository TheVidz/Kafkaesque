package memory

import (
	"fmt"
	"github.com/TheVidz/Kafkaesque/pkg/broker"
	"sync"
	"time"
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
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	t.Broadcast(msg)
	return nil
}

func (mb *MemoryBroker) Subscribe(topic string, ch chan<- broker.Message) error {
	t := mb.getOrCreateTopic(topic)
	t.AddSubscriber(ch)
	return nil
}
