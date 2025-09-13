package main

import (
	"fmt"
	"time"

	"github.com/TheVidz/Kafkaesque/pkg/broker"
	"github.com/TheVidz/Kafkaesque/pkg/broker/memory"
)
func main() {
	b := memory.NewMemoryBroker()

	// === Consumers ===
	// Consumer A subscribes to "orders"
	subA := make(chan broker.Message, 10)
	b.Subscribe("orders", subA)
	go consumer("Consumer-A", subA)

	// Consumer B subscribes to "orders"
	subB := make(chan broker.Message, 10)
	b.Subscribe("orders", subB)
	go consumer("Consumer-B", subB)

	// Consumer C subscribes to "payments"
	subC := make(chan broker.Message, 10)
	b.Subscribe("payments", subC)
	go consumer("Consumer-C", subC)

	// === Producers ===
	// Producer 1 publishes to "orders"
	go func() {
		for i := 1; i <= 3; i++ {
			msg := fmt.Sprintf("order-%d", i)
			b.Publish("orders", []byte(msg))
			fmt.Printf("[Producer-1] published: %s\n", msg)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Producer 2 publishes to "payments"
	go func() {
		for i := 1; i <= 2; i++ {
			msg := fmt.Sprintf("payment-%d", i)
			b.Publish("payments", []byte(msg))
			fmt.Printf("[Producer-2] published: %s\n", msg)
			time.Sleep(700 * time.Millisecond)
		}
	}()

	// Give enough time for goroutines to finish
	time.Sleep(3 * time.Second)
}

func consumer(name string, ch <-chan broker.Message) {
	for msg := range ch {
		fmt.Printf("[%s] received from topic=%s : %s\n",
			name, msg.Topic, string(msg.Payload))
	}
}