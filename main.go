package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

const (
	kafkaTopic  = "example-topic"
	kafkaBroker = "localhost:9092"
)

func produce(ctx context.Context, message string) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", time.Now().Unix())),
			Value: []byte(message),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaBroker},
		Topic:     kafkaTopic,
		GroupID:   "example-group",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	if message == "" {
		http.Error(w, "Please provide a message", http.StatusBadRequest)
		return
	}

	err := produce(context.Background(), message)
	if err != nil {
		http.Error(w, "Failed to produce message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Produced message: %s\n", message)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/send", messageHandler).Methods("GET")

	go func() {
		log.Println("Starting consumer...")
		consume(context.Background())
	}()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// import (
// 	"bytes"
// 	"encoding/base64"
// 	"encoding/binary"
// 	"fmt"
// )

// // KafkaMessage represents a simplified Kafka message
// type KafkaMessage struct {
// 	Size        int32
// 	MagicByte   byte
// 	Attributes  byte
// 	Timestamp   int64
// 	KeyLength   int32
// 	Key         []byte
// 	ValueLength int32
// 	Value       []byte
// }

// // DecodeKafkaMessage decodes a byte slice into a KafkaMessage
// func DecodeKafkaMessage(data []byte) (*KafkaMessage, error) {
// 	buffer := bytes.NewBuffer(data)

// 	var message KafkaMessage
// 	if err := binary.Read(buffer, binary.BigEndian, &message.Size); err != nil {
// 		return nil, err
// 	}

// 	message.MagicByte, _ = buffer.ReadByte()
// 	message.Attributes, _ = buffer.ReadByte()

// 	if err := binary.Read(buffer, binary.BigEndian, &message.Timestamp); err != nil {
// 		return nil, err
// 	}

// 	if err := binary.Read(buffer, binary.BigEndian, &message.KeyLength); err != nil {
// 		return nil, err
// 	}

// 	if message.KeyLength > 0 {
// 		message.Key = make([]byte, message.KeyLength)
// 		if _, err := buffer.Read(message.Key); err != nil {
// 			return nil, err
// 		}
// 	}

// 	if err := binary.Read(buffer, binary.BigEndian, &message.ValueLength); err != nil {
// 		return nil, err
// 	}

// 	if message.ValueLength > 0 {
// 		message.Value = make([]byte, message.ValueLength)
// 		if _, err := buffer.Read(message.Value); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return &message, nil
// }

// func main() {

// 	data, err := base64.StdEncoding.DecodeString("AAAAMQAAAAEAAAABAA1leGFtcGxlLXRvcGljAAAAAQAAAAAAAP//////////AAAAAAAAAAA=")
// 	if err != nil {
// 		fmt.Println("failed to decode the data", err)
// 	}
// 	// Example binary data (not a real Kafka message)

// 	message, err := DecodeKafkaMessage(data)
// 	if err != nil {
// 		fmt.Println("Error decoding message:", err)
// 		return
// 	}

// 	fmt.Printf("Decoded message: %+v\n", message)
// }
