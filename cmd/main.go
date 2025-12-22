package main

import (
	"Trade-app/internal/api"
	"Trade-app/internal/engine"
	"Trade-app/internal/kafka"
	"Trade-app/internal/redis"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	book := engine.NewOrderBook()
	//kafka producer
	kafkaProducer := kafka.NewProducer([]string{"localhost:9092"}, "orders")
	// cache for redis
	cache := redis.NewCache("localhost:6379")

	handler := &api.Handler{
		Book:  book,
		Kafka: kafkaProducer,
		Redis: cache,
		Ctx:   ctx,
	}

	router := api.NewRouter(handler)

	log.Println("Server running at :8080")
	router.Run(":8080")
}
