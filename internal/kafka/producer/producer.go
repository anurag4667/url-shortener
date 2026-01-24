package producer

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type ClickProducer struct {
	writer *kafka.Writer
}

func NewClickProducer(brokers []string) *ClickProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        "clicks",
		Balancer:     &kafka.Hash{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Async:        true,
	})

	return &ClickProducer{writer: writer}
}

func (p *ClickProducer) Publish(ctx context.Context, shortCode string, url string) {
	msg := kafka.Message{
		Key:   []byte(shortCode),
		Value: []byte(url),
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("failed to publish click event: %v", err)
	}

	log.Printf("published click event for shortCode: %s", shortCode)
}

func (p *ClickProducer) Close() error {
	return p.writer.Close()
}
