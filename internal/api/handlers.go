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

	if err := h.Kafka.Publish(h.Ctx, order.ID, []byte("order executed")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kafka publish failed"})
		return
	}

	if err := h.Redis.Set(h.Ctx, "orderbook", []byte("snapshot")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis save failed"})
		return
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
