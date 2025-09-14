package broker

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Broker struct {
	mu     sync.RWMutex
	topics map[string]*Topic
	seq    int64
}

// NewBroker constructs the broker.
func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
		seq:    0,
	}
}

// CreateTopic ensures topic exists.
func (b *Broker) CreateTopic(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.topics[name]; !ok {
		b.topics[name] = NewTopic(name)
	}
}

// Subscribe registers a Subscriber to a topic (creates topic if needed).
func (b *Broker) Subscribe(topic string, sub *Subscriber) error {
	if sub == nil {
		return errors.New("subscriber is nil")
	}
	b.CreateTopic(topic)
	b.mu.RLock()
	t := b.topics[topic]
	b.mu.RUnlock()
	t.AddSubscriber(sub)
	return nil
}

// Publish adds a message and starts delivery attempts.
func (b *Broker) Publish(topic string, payload []byte) (Message, error) {
	b.mu.Lock()
	b.seq++
	msg := Message{
		ID:        b.seq,
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now(),
	}
	b.mu.Unlock()

	b.mu.RLock()
	t, ok := b.topics[topic]
	b.mu.RUnlock()
	if !ok {
		return Message{}, fmt.Errorf("topic %s not found", topic)
	}

	t.Publish(msg)
	return msg, nil
}

// Convenience: publish synchronously and print msg ID
func (b *Broker) PublishString(topic, body string) (Message, error) {
	return b.Publish(topic, []byte(body))
}
