package broker

type Subscriber struct {
	ID    string
	Ch    chan Message // broker -> subscriber
	ackCh chan int64   // subscriber -> broker (acks)
}

func NewSubscriber(id string) *Subscriber {
	return &Subscriber{
		ID:    id,
		Ch:    make(chan Message, 16),
		ackCh: make(chan int64, 16),
	}
}

// Call this from subscriber code after processing a message
func (s *Subscriber) Ack(msgID int64) {
	s.ackCh <- msgID
}
