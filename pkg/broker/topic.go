package broker

import (
	"fmt"
	"sync"
	"time"
)

type Topic struct {
	name          string
	mu            sync.Mutex
	subscribers   map[string]*Subscriber
	pending       map[int64]Message         // pending messages
	ackTracker    map[int64]map[string]bool // msgID -> subID -> acked
	retryInterval time.Duration
}

func NewTopic(name string) *Topic {
	return &Topic{
		name:          name,
		subscribers:   make(map[string]*Subscriber),
		pending:       make(map[int64]Message),
		ackTracker:    make(map[int64]map[string]bool),
		retryInterval: 2 * time.Second,
	}
}

func (t *Topic) AddSubscriber(sub *Subscriber) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.subscribers[sub.ID] = sub
}

func (t *Topic) Publish(msg Message) {
	t.mu.Lock()
	// register as pending and prepare ack tracker
	t.pending[msg.ID] = msg
	t.ackTracker[msg.ID] = make(map[string]bool)
	for sid := range t.subscribers {
		t.ackTracker[msg.ID][sid] = false
	}
	// copy subscribers slice to avoid holding lock for long
	subs := make([]*Subscriber, 0, len(t.subscribers))
	for _, s := range t.subscribers {
		subs = append(subs, s)
	}
	t.mu.Unlock()

	// deliver to each subscriber in own goroutine (retries until ack)
	for _, s := range subs {
		go t.deliverWithRetry(s, msg)
	}
}

func (t *Topic) deliverWithRetry(sub *Subscriber, msg Message) {
	for {
		select {
		case sub.Ch <- msg:
			// wait for ack or timeout
			select {
			case ackID := <-sub.ackCh:
				if ackID != msg.ID {
					// ignore unrelated ack
					continue
				}
				// mark acked for this subscriber
				t.mu.Lock()
				t.ackTracker[msg.ID][sub.ID] = true
				// check if all acked
				allAcked := true
				for sid := range t.ackTracker[msg.ID] {
					if !t.ackTracker[msg.ID][sid] {
						allAcked = false
						break
					}
				}
				if allAcked {
					delete(t.pending, msg.ID)
					delete(t.ackTracker, msg.ID)
					// optionally log or notify
					fmt.Printf("[Topic %s] message %d fully acked by all subscribers\n", t.name, msg.ID)
				}
				t.mu.Unlock()
				return
			case <-time.After(t.retryInterval):
				// timeout waiting for ack -> retry send
				fmt.Printf("[Topic %s] timeout waiting ack for msg %d to sub %s, retrying\n", t.name, msg.ID, sub.ID)
				continue
			}
		case <-time.After(t.retryInterval):
			// subscriber not ready to receive -> retry
			fmt.Printf("[Topic %s] subscriber %s not ready for msg %d, retrying\n", t.name, sub.ID, msg.ID)
			continue
		}
	}
}
