package kafka

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}

func (c *Consumer) Consume(ctx *gin.Context) {
	msg, err := c.reader.ReadMessage(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read kafka message",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Sucess":"new event from kafka",
		"topic": msg.Topic,
		"key":   string(msg.Key),
		"value": string(msg.Value),
	})
}
