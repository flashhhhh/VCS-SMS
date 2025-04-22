package main

import (
	"fake_server_service/infrastructure/redis"
	"fake_server_service/internal/repository"
	"fake_server_service/internal/service"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
)

var MAX_SERVERS = 10

func main() {
	// Initialize the logger for the fake_server_service
	currentPath, _ := os.Getwd()
	fakeServerServiceLogPath := filepath.Join(currentPath, "logs", "fake_server_service.log")
	logging.InitLogger("fake_server_service", fakeServerServiceLogPath, 10, 5, 30)

	// Load running environment variable
	environment := env.GetEnv("RUNNING_ENVIRONMENT", "local")
	logging.LogMessage("fake_server_service", "Running in "+environment+" environment", "INFO")

	// Load environment variables from the .env file
	environmentFilePath := filepath.Join(currentPath, "configs", environment+".env")
	if err := env.LoadEnv(environmentFilePath); err != nil {
		logging.LogMessage("fake_server_service", "Failed to load environment variables from "+environmentFilePath+": "+err.Error(), "FATAL")
		logging.LogMessage("fake_server_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("fake_server_service", "Environment variables loaded successfully from "+environmentFilePath, "INFO")
	}

	// Connect to redis
	redisAddress := env.GetEnv("REDIS_HOST", "localhost") + ":" + env.GetEnv("REDIS_PORT", "6379")
	redis := redis.NewRedisClient(redisAddress)

	// Initialize the fake server service
	fakeServerRepository := repository.NewFakeServerRepository(redis)
	fakeServerService := service.NewFakeServerService(&fakeServerRepository)

	// Delete all existing servers
	if err := fakeServerService.DeleteServers(); err != nil {
		logging.LogMessage("fake_server_service", "Failed to delete all servers: "+err.Error(), "ERROR")
		logging.LogMessage("fake_server_service", "Exiting the program...", "ERROR")
		os.Exit(1)
	} else {
		logging.LogMessage("fake_server_service", "All servers deleted successfully", "INFO")
	}

	for {
		runningServers, _ := fakeServerService.CountRunningServers()
		logging.LogMessage("fake_server_service", "Number of running servers: "+strconv.Itoa(runningServers), "INFO")

		for i := 0; i < MAX_SERVERS-runningServers; i++ {
			// Randomly select a server ID
			var serverID int
			for {
				serverID = rand.Intn(65536) // Random number between 0 and 65535
				logging.LogMessage("fake_server_service", "Selected server ID: "+strconv.Itoa(serverID), "DEBUG")

				// Check if the server can be enabled
				canBeEnabled, _ := fakeServerService.CheckServer(serverID)
				if !canBeEnabled {
					logging.LogMessage("fake_server_service", "Server with ID "+strconv.Itoa(serverID)+" can't be enabled", "INFO")
					continue
				}

				break
			}

			// Enable the server
			go func () {
				fakeServerService.HostServer(serverID, 2000)	
			}()

			logging.LogMessage("fake_server_service", "Server with ID "+strconv.Itoa(serverID)+" has been enabled", "INFO")
			logging.LogMessage("fake_server_service", "Waiting 10 seconds before checking again...", "INFO")
		}

		// Sleep for 60 seconds before checking again
		logging.LogMessage("fake_server_service", "Sleeping for 60 seconds before checking again...", "INFO")
		time.Sleep(60 * time.Second)
	}
}