package main

import (
	"context"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "example4",
		Partition: 0,
		MinBytes:  10,
		MaxBytes:  10e6,
	})
	defer r.Close()

	m, err := r.ReadMessage(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message offset: %d, value: %s\n", m.Offset, string(m.Value))

	// Simplified decoding for demonstration
	fmt.Println("Decoding message value:")
	if len(m.Value) >= 2 {
		magicByte := m.Value[0]
		attributes := m.Value[1]
		fmt.Printf("Magic Byte: %d, Attributes: %d\n", magicByte, attributes)
		if len(m.Value) > 2 {
			fmt.Printf("Message payload: %s\n", string(m.Value))
		}
	}
	time.Sleep(5 * time.Second)
}
