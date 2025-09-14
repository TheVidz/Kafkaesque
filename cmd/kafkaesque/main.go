package main

import (
	"fmt"
	"time"

	"github.com/TheVidz/Kafkaesque/pkg/broker"
)
func main() {
	b := broker.NewBroker()
	b.CreateTopic("demo")

	// subscriber that ACKs after a delay (shows retry behavior)
	sub1 := broker.NewSubscriber("sub1")
	_ = b.Subscribe("demo", sub1)

	// subscriber that immediately ACKs (fast)
	sub2 := broker.NewSubscriber("sub2")
	_ = b.Subscribe("demo", sub2)

	// sub1 goroutine: simulate slow processing then ack
	go func() {
		for msg := range sub1.Ch {
			fmt.Printf("[sub1] received msg %d: %s\n", msg.ID, string(msg.Payload))
			// simulate slow processing so broker will retry
			time.Sleep(3 * time.Second)
			fmt.Printf("[sub1] ACKing %d\n", msg.ID)
			sub1.Ack(msg.ID)
		}
	}()

	// sub2 goroutine: fast processing
	go func() {
		for msg := range sub2.Ch {
			fmt.Printf("[sub2] received msg %d: %s\n", msg.ID, string(msg.Payload))
			// immediate ack
			sub2.Ack(msg.ID)
		}
	}()

	// publish a few messages
	for i := 1; i <= 2; i++ {
		body := fmt.Sprintf("message-%d", i)
		msg, err := b.PublishString("demo", body)
		if err != nil {
			fmt.Printf("publish err: %v\n", err)
			continue
		}
		fmt.Printf("[publisher] published id=%d body=%s\n", msg.ID, body)
	}

	// wait so we can observe retries
	time.Sleep(12 * time.Second)
	fmt.Println("demo finished")
}