package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) *redis.Client {
	log.Printf("Connecting to Redis at %s", addr)
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}