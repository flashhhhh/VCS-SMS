package main

import (
	"net/http"
	"os"
	"path/filepath"
	"user_service/infrastructure/postgres"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
)

func main() {
	// Initialize the logger for the user_service
	currentPath, _ := os.Getwd()
	userServiceLogPath := filepath.Join(currentPath, "logs", "user_service.log")
	logging.InitLogger("user_service", userServiceLogPath, 10, 5, 30)

	// Load running environment variable
	environment := env.GetEnv("RUNNING_ENVIRONMENT", "local")
	logging.LogMessage("user_service", "Running in "+environment+" environment", "INFO")

	// Load environment variables from the .env file
	environmentFilePath := filepath.Join(currentPath, "configs", environment+".env")
	if err := env.LoadEnv(environmentFilePath); err != nil {
		logging.LogMessage("user_service", "Failed to load environment variables from "+environmentFilePath+": "+err.Error(), "FATAL")
		logging.LogMessage("user_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("user_service", "Environment variables loaded successfully from "+environmentFilePath, "INFO")
	}

	// Connect to the database
	dsn := "host=" + env.GetEnv("POSTGRES_HOST", "localhost") +
		" user=" + env.GetEnv("POSTGRES_USER", "postgres") +
		" password=" + env.GetEnv("POSTGRES_PASSWORD", "password") +
		" dbname=" + env.GetEnv("POSTGRES_NAME", "user_service") +
		" port=" + env.GetEnv("POSTGRES_PORT", "5432") +
		" sslmode=disable"
	db := postgres.ConnectDB(dsn)

	// Migrate the database
	postgres.Migrate(db)

	// Start the HTTP server
	user_service_port := env.GetEnv("USER_SERVICE_PORT", "10001")

	logging.LogMessage("user_service", "Starting HTTP server on port " + user_service_port + "...", "INFO")
	http.ListenAndServe(":" + user_service_port, nil)
}