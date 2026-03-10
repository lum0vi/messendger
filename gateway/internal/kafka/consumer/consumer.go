package consumer

import (
	"context"
	"encoding/json"
	"time"

	"gateway/internal/handler"
	"gateway/internal/models"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

const (
	TopicMessageDelivered = "message.delivered"
)

type ConsumerKafka struct {
	reader *kafka.Reader
	h      *handler.Handler
}

func NewConsumerMessage(brokers []string, h *handler.Handler) *ConsumerKafka {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       TopicMessageDelivered,
		GroupID:     "gateway-group",
		StartOffset: kafka.FirstOffset,
		MinBytes:    1e3,  // 1KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
	})
	return &ConsumerKafka{
		reader: r,
		h:      h,
	}
}

func (c *ConsumerKafka) Start(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			logrus.Errorf("Error reading message: %v", err)
			continue
		}
		logrus.Infof("Message received from topic %s: %s", m.Topic, string(m.Value))
		var req models.MessageDelivery
		if err = json.Unmarshal(m.Value, &req); err != nil {
			logrus.Errorf("Error unmarshalling message: %v", err)
		}

		if err := c.h.HandleWebsocketRequest(&req); err != nil {
			logrus.Errorf("Error handling message: %v", err)
			continue
		}
		logrus.Infof("Message delivered from topic %s: %+v", m.Topic, req)
	}
}

func (c *ConsumerKafka) Close() error {
	return c.reader.Close()
}
