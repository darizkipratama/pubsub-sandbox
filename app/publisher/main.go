package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "contoh-topic",
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	for i := 0; i < 10; i++ {
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(fmt.Sprintf("Key-%d", i)),
				Value: []byte(fmt.Sprintf("Pesan ke-%d", i)),
			},
		)
		if err != nil {
			fmt.Println("Gagal kirim pesan:", err)
		} else {
			fmt.Println("Berhasil kirim pesan ke Kafka:", i)
		}
		time.Sleep(time.Second)
	}

	// Mengirim Pesan Yang Berisi Pesan error, untuk simulasi gagal proses
	for i := 0; i < 2; i++ {
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(fmt.Sprintf("Key-%d", i)),
				Value: []byte(fmt.Sprintf("Pesan error ke-%d", i)),
			},
		)
		if err != nil {
			fmt.Println("Gagal kirim pesan:", err)
		} else {
			fmt.Println("Berhasil kirim pesan ke Kafka:", i)
		}
		time.Sleep(time.Second)
	}
}
