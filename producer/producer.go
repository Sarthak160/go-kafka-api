package main

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "example4",
	})
	defer w.Close()

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte("Hello Kafka 4"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Message produced")
}
