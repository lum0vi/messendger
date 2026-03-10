package producer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	TopicMessageSend = "message.send"
)

type ProducerKafka struct {
	writer *kafka.Writer
}

func NewProducerMessage(brokers []string) *ProducerKafka {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.Hash{},
		Topic:    TopicMessageSend,
	}
	return &ProducerKafka{
		writer: writer,
	}
}

func (p *ProducerKafka) Send(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})
}

func (p *ProducerKafka) Close() error {
	return p.writer.Close()
}
