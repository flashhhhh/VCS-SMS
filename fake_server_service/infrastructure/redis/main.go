package redis

import (
	"context"
	"os"

	"github.com/flashhhhh/pkg/logging"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) *redis.Client {
	logging.LogMessage("fake_server_service", "Connecting to Redis...", "INFO")
	logging.LogMessage("fake_server_service", "Redis Address: "+addr, "DEBUG")

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		logging.LogMessage("fake_server_service", "Failed to connect to Redis: "+err.Error(), "FATAL")
		logging.LogMessage("fake_server_service", "Exiting the program...", "FATAL")
		
		os.Exit(1)
	} else {
		logging.LogMessage("fake_server_service", "Connected to Redis successfully", "INFO")
	}

	logging.LogMessage("fake_server_service", "Redis client created successfully", "INFO")
	return client
}