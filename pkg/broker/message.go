package broker

import "time"

type Message struct {
	ID        int64
	Topic     string
	Payload   []byte
	Timestamp time.Time
}
