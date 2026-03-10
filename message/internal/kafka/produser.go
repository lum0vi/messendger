package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type ProducerMessage struct {
	writer *kafka.Writer
}

func NewProducerMessage(brokers []string) *ProducerMessage {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.Hash{},
		Topic:    TopicMessageDelivered,
	}
	return &ProducerMessage{
		writer: writer,
	}
}

func (p *ProducerMessage) Send(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})
}
func (p *ProducerMessage) Close() error {
	return p.writer.Close()
}
