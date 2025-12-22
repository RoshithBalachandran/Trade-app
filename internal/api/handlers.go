package api

import (
	"Trade-app/internal/candle"
	"Trade-app/internal/engine"
	"Trade-app/internal/kafka"
	"Trade-app/internal/redis"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Book      *engine.OrderBook
	Kafka     *kafka.Producer
	Redis     *redis.Cache
	Ctx       context.Context
	CandleAgg *candle.OHLCVAggregator
}

func (h *Handler) PlaceOrder(c *gin.Context) {
	var order engine.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trades := h.Book.Match(&order)

	// Publish to Kafka
	err := h.Kafka.Publish(h.Ctx, order.ID, []byte("order executed"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to send event to kafka"})
	}

	// Update Redis snapshot
	err = h.Redis.Set(h.Ctx, "orderbook", []byte("snapshot"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to save order book to redis"})
	}

	c.JSON(http.StatusOK, trades)
}

func (h *Handler) GetOrderBook(c *gin.Context) {
	c.JSON(http.StatusOK, h.Book)
}

func (h *Handler) GetOHLCVCandles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"1s": h.CandleAgg.OneSecCandles,
		"1m": h.CandleAgg.OneMinCandles,
	})
}
