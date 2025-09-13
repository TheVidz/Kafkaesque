package broker

type Broker interface {
	Publish(topic string, payload []byte) error
	Subscribe(topic string, ch chan<- Message) error
}
