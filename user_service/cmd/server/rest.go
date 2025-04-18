package main

import (
	"net/http"
	"os"
	"path/filepath"
	"user_service/infrastructure/postgres"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
	"gorm.io/gorm"
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
	
	dbChan := make(chan *gorm.DB)

	go func() {
		db := postgres.ConnectDB(dsn)
		dbChan <- db
	}()
	db := <-dbChan

	// Run migrations
	postgres.Migrate(db)

	// Start the HTTP server
	logging.LogMessage("user_service", "Starting HTTP server on port 8080...", "INFO")
	http.ListenAndServe(":8080", nil)
}