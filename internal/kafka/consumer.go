package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			StartOffset: kafka.FirstOffset,
			MinBytes:    1e3,
			MaxBytes:    10e6,
		}),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopping")
			_ = c.reader.Close()
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Kafka fetch error: %v", err)
				continue
			}

			log.Printf(
				"Consumed event | topic=%s key=%s value=%s",
				msg.Topic,
				string(msg.Key),
				string(msg.Value),
			)

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Kafka commit error: %v", err)
			}
		}
	}
}
