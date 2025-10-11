package main

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "contoh-topic",
		GroupID: "group-darizki",
	})

	defer reader.Close()

	fmt.Println("Menunggu pesan dari Kafka...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Gagal baca pesan:", err)
			break
		}
		fmt.Printf("Pesan diterima: %s | Key: %s\n", string(msg.Value), string(msg.Key))
	}
}
