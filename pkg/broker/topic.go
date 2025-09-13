package broker

import (
	"sync"
)

type Topic struct {
	Name        string
	subscribers []chan<- Message
	mu          sync.RWMutex
}

func NewTopic(name string) *Topic {
	return &Topic{
		Name:        name,
		subscribers: make([]chan<- Message, 0),
	}
}

func (t *Topic) AddSubscriber(ch chan<- Message) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.subscribers = append(t.subscribers, ch)
}

func (t *Topic) Broadcast(msg Message) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, sub := range t.subscribers {
		// non-blocking send
		select {
		case sub <- msg:
		default:
			// drop if subscriber is slow
		}
	}
}
