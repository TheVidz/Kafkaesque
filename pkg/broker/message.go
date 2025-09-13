package broker

import "time"

type Message struct {
	ID        string
	Topic     string
	Payload   []byte
	Timestamp time.Time
}
