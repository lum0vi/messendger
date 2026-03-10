package kafka

import (
	"context"
	"message/internal/repository"
)

type ConsumerMess interface {
	Start(ctx context.Context)
	Close() error
}
type ProducerMess interface {
	Send(ctx context.Context, key, value []byte) error
	Close() error
}

type ConsumerKafka struct {
	Consumer ConsumerMess
}

type ProducerKafka struct {
	Producer ProducerMess
}

func NewConsumerKafka(addr []string, repo *repository.Repository, prod *ProducerKafka) *ConsumerKafka {
	return &ConsumerKafka{
		Consumer: NewConsumerMessage(addr, repo, prod),
	}
}

func NewProducerKafka(addr []string) *ProducerKafka {
	return &ProducerKafka{
		Producer: NewProducerMessage(addr),
	}
}
