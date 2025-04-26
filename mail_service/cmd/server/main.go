package main

import (
	"mail_service/internal/service"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
)

func main() {
	// Initialize logger for mail_service
	currentPath, _ := os.Getwd()
	mailServiceLogPath := filepath.Join(currentPath, "logs", "mail_service.log")
	logging.InitLogger("mail_service", mailServiceLogPath, 10, 5, 30)

	// Load running environment variable
	environment := env.GetEnv("RUNNING_ENVIRONMENT", "local")
	logging.LogMessage("mail_service", "Running in "+environment+" environment", "INFO")

	// Load environment variables from the .env file
	environmentFilePath := filepath.Join(currentPath, "configs", environment+".env")
	if err := env.LoadEnv(environmentFilePath); err != nil {
		logging.LogMessage("mail_service", "Failed to load environment variables from "+environmentFilePath+": "+err.Error(), "FATAL")
		logging.LogMessage("mail_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("mail_service", "Environment variables loaded successfully from "+environmentFilePath, "INFO")
	}

	go func () {
		numServers := 10
		numOnServers := 7
		numOffServers := 3
		meanUptimeRate := 70.0

		to := env.GetEnv("SERVER_ADMINISTRATOR_EMAIL", "")
		subject := "Daily Server Status Report for " + time.Now().Format("2006-01-02")

		service.PrepareEmail(to, subject, numServers, numOnServers, numOffServers, meanUptimeRate)

		time.Sleep(24 * time.Hour)
	}()

	// Listen for interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs // Wait for interrupt
}