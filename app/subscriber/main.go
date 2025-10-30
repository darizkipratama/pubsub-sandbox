package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	maxRetries = 3
	dlqTopic   = "contoh-topic-dlq"
)

// fungsi yang mensimulasikan pemrosesan pesan dari publisher
func processMessage(msg kafka.Message) error {
	// Untuk keperluan simulasi, kita anggap pesan yang mengandung kata 'error' akan gagal diproses
	if containsError := strings.Contains(strings.ToLower(string(msg.Value)), "error"); containsError {
		return fmt.Errorf("failed to process message containing 'error': %s", string(msg.Value))
	}
	return nil
}

// fungsi yang mengimplementasikan mekanisme retry dalam pemrosesan pesan
func processMessageWithRetry(msg kafka.Message, maxRetries int) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = processMessage(msg)
		if err == nil {
			return nil
		}

		fmt.Printf("Percobaan Pemrosesan Pesan %d Gagal: %v\n", attempt, err)
		if attempt < maxRetries {
			// Wait before retrying, with exponential backoff
			time.Sleep(time.Duration(attempt*attempt) * time.Second)
		}
	}
	return err
}

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "contoh-topic",
		GroupID: "group-darizki",
	})

	// Create a writer for the Dead Letter Queue
	dlqWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   dlqTopic,
	})

	defer reader.Close()
	defer dlqWriter.Close()

	fmt.Println("Menunggu pesan dari Kafka...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Gagal baca pesan:", err)
			break
		}

		// 1. Coba proses pesan dengan mekanisme retry
		err = processMessageWithRetry(msg, maxRetries)
		if err != nil {
			fmt.Printf("Jika pesan gagal diproses setelah %d percobaan, akan disimpan ke DLQ\n", maxRetries)

			// Menambahkan header informasi error
			headers := append(msg.Headers,
				kafka.Header{Key: "error", Value: []byte(err.Error())},
				kafka.Header{Key: "original_topic", Value: []byte(msg.Topic)},
				kafka.Header{Key: "failed_at", Value: []byte(time.Now().String())},
			)

			// 2. Jika tetap Gagal maka Kirim ke DLQ
			err = dlqWriter.WriteMessages(context.Background(), kafka.Message{
				Key:     msg.Key,
				Value:   msg.Value,
				Headers: headers,
			})
			if err != nil {
				fmt.Printf("Gagal mengirim pesan ke topic DLQ: %v\n", err)
			} else {
				fmt.Printf("Pesan dikirim ke topic DLQ: %s\n", string(msg.Value))
			}
			continue
		}

		fmt.Printf("Pesan berhasil diproses: %s | Key: %s\n", string(msg.Value), string(msg.Key))
	}
}
