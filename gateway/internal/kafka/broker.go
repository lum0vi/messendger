package kafka

//import (
//	"context"
//	"/gateway/internal/handler"
//	"/gateway/internal/kafka/consumer"
//)
//
//type ConsumerInterface interface {
//	Start(context.Context)
//	Close() error
//}
//type Consumer struct {
//	Consumer ConsumerInterface
//}
//
//type ProducerInterface interface {
//	Send(ctx context.Context, key, value []byte) error
//	Close() error
//}
//type Producer struct {
//	Producer ProducerInterface
//}
//
//func NewConsumer(addr []string, h *handler.Handler) *Consumer {
//	return &Consumer{
//		Consumer: consumer.NewConsumerMessage(addr, h),
//	}
//}
//func NewProducer(addr []string) *Producer {
//	return &Producer{
//		Producer: NewProducerMessage(addr),
//	}
//}
