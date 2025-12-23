package api

import (
	"Trade-app/internal/kafka"
	"Trade-app/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handler, Consumer kafka.Consumer) *gin.Engine {
	r := gin.Default()
	// Public routes
	r.POST("/signin", Signin)
	r.POST("/login", Login)
	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/order", handler.PlaceOrder)
		auth.GET("/orderbook", handler.GetOrderBook)
		auth.GET("/candles", handler.GetOHLCVCandles)

		auth.GET("/consumer", Consumer.Consume)
	}
	return r
}
