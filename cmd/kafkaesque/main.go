package main

import (
	"fmt"
	"time"

	"github.com/TheVidz/Kafkaesque/pkg/broker"
	"github.com/TheVidz/Kafkaesque/pkg/broker/memory"
)

func main() {
	b := memory.NewMemoryBroker()

	// Create subscriber
	subCh := make(chan broker.Message, 10)
	b.Subscribe("orders", subCh)

	// Consumer goroutine
	go func() {
		for msg := range subCh {
			fmt.Printf("[Consumer] got message: %s\n", msg.Payload)
		}
	}()

	// Producer loop
	for i := 0; i < 5; i++ {
		b.Publish("orders", []byte(fmt.Sprintf("order-%d", i)))
		time.Sleep(500 * time.Millisecond)
	}
}
