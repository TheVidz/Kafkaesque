package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/TheVidz/Kafkaesque/pkg/broker"
// 	"github.com/TheVidz/Kafkaesque/pkg/broker/memory"
// )

// func main() {
//     b := memory.New()
//     _ = b.CreateTopic(context.Background(), "demo", broker.TopicOptions{Partitions: 1})

//     msgs, stop, err := b.Consumer().Subscribe(context.Background(), "demo", "g1")
//     if err != nil { panic(err) }
//     defer stop()

//     go func() {
//         for m := range msgs {
//             fmt.Printf("[consumer] %s:%d %s\n", m.Topic, m.Offset, string(m.Value))
//         }
//     }()

//     for i := 0; i < 3; i++ {
//         _, _ = b.Producer().Publish(context.Background(), "demo", "", []byte(fmt.Sprintf("hello-%d", i)))
//         time.Sleep(200 * time.Millisecond)
//     }
//     time.Sleep(time.Second)
// }
