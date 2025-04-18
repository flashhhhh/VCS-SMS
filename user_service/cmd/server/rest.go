package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
	"user_service/infrastructure/postgres"

	"github.com/flashhhhh/pkg/logging"
)

func main() {
	// Initialize the logger for the user_service
	currentPath, _ := os.Getwd()
	userServiceLogPath := filepath.Join(currentPath, "logs", "user_service.log")
	logging.InitLogger("user_service", userServiceLogPath, 10, 5, 30)

	// Connect to the database
	dsn := "host=localhost user=postgres password=yourpassword dbname=user_service port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := postgres.ConnectDB(dsn)

	if err != nil {
		go func() {
			for {
				logging.LogMessage("user_service", "Trying to connect to the database...", "INFO")
				db, err = postgres.ConnectDB(dsn)

				if err == nil {
					logging.LogMessage("user_service", "Connected to the database successfully", "INFO")
					break
				}

				logging.LogMessage("user_service", "Failed to connect to the database: "+err.Error(), "ERROR")
				logging.LogMessage("user_service", "Retrying in 10 seconds...", "INFO")
				time.Sleep(10 * time.Second)
			}
		}()
	} else {
		logging.LogMessage("user_service", "Connected to the database successfully", "INFO")
	}

	// Run migrations
	go func () {
		for {
			if db == nil {
				logging.LogMessage("user_service", "Database connection is nil, retrying migration in 10 seconds...", "ERROR")
				time.Sleep(10 * time.Second)
				continue
			}

			err = postgres.Migrate(db)
			if err != nil {
				// Fatal error, exit the program
				logging.LogMessage("user_service", "Failed to run migrations: "+err.Error(), "ERROR")
				logging.LogMessage("user_service", "Exiting the program...", "FATAL")
				os.Exit(1)
			}

			logging.LogMessage("user_service", "Migrations completed successfully", "INFO")
			break
		}
	}()

	// Start the HTTP server
	logging.LogMessage("user_service", "Starting HTTP server on port 8080...", "INFO")
	http.ListenAndServe(":8080", nil)
}