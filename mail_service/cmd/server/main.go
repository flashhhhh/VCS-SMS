package main

import (
	"mail_service/api/routes"
	"mail_service/infrastructure/grpc"
	grpcclient "mail_service/internal/grpc_client"
	"mail_service/internal/handler"
	"mail_service/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/logging"
	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
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

	// Initialize gRPC client
	grpcClient, err := grpc.StartGRPCClient()
	if err != nil {
		panic(err)
	}

	client := grpcclient.NewServerAdministrationServiceClient(grpcClient)
	resp, err := client.GetServerInformation(
		time.Now().Add(-24*time.Hour).Unix(),
		time.Now().Unix(),
	)

	if err != nil {
		println("Failed to get server information:", err.Error())
		panic(err)
	}

	println("Server Information:", resp)

	mailService := service.NewMailService(client)
	mailHandler := handler.NewMailHandler(mailService)

	// Send email report every 24 hours
	go func () {
		mailService.StartEmailReport(time.Now().Add(-24*time.Hour).Unix(), time.Now().Unix())
		time.Sleep(24 * time.Hour)
	}()

	// Start the server
	serverHost := env.GetEnv("MAIL_SERVICE_HOST", "localhost")
	serverPort := env.GetEnv("MAIL_SERVICE_PORT", "10003")

	r := mux.NewRouter()
	routes.RegisterRoutes(r, mailHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, change this for security
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	logging.LogMessage("mail_service", "Starting server on "+serverHost+":"+serverPort, "INFO")
	if err := http.ListenAndServe(serverHost+":"+serverPort, corsHandler); err != nil {
		logging.LogMessage("mail_service", "Failed to start server: "+err.Error(), "FATAL")
		logging.LogMessage("mail_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}

	logging.LogMessage("mail_service", "Server stopped", "INFO")
	logging.LogMessage("mail_service", "Exiting the program...", "INFO")
	os.Exit(0)
}