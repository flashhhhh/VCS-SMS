package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"server_administration_service/infrastructure/grpc"
	"server_administration_service/infrastructure/postgres"
	"server_administration_service/infrastructure/redis"
	"server_administration_service/internal/handler"
	"server_administration_service/internal/repository"
	"server_administration_service/internal/service"
	"syscall"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/kafka"
	"github.com/flashhhhh/pkg/logging"
)

func main() {
	// Initialize the logger for the server_administration_service
	currentPath, _ := os.Getwd()
	serverAdministrationServiceLogPath := filepath.Join(currentPath, "logs", "server_administration_service.log")
	logging.InitLogger("server_administration_service", serverAdministrationServiceLogPath, 10, 5, 30)

	// Load running environment variable
	environment := env.GetEnv("ENVIRONMENT", "local")
	logging.LogMessage("server_administration_service", "Running in "+environment+" environment", "INFO")

	// Load environment variables from the .env file
	environmentFilePath := filepath.Join(currentPath, "configs", environment+".env")
	if err := env.LoadEnv(environmentFilePath); err != nil {
		logging.LogMessage("server_administration_service", "Failed to load environment variables from "+environmentFilePath+": "+err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	} else {
		logging.LogMessage("server_administration_service", "Environment variables loaded successfully from "+environmentFilePath, "INFO")
	}

	// Connect to the database
	dsn := "host=" + env.GetEnv("SERVER_POSTGRES_HOST", "localhost") +
		" port=" + env.GetEnv("SERVER_POSTGRES_PORT", "5432") +
		" user=" + env.GetEnv("SERVER_POSTGRES_USER", "postgres") +
		" password=" + env.GetEnv("SERVER_POSTGRES_PASSWORD", "password") +
		" dbname=" + env.GetEnv("SERVER_POSTGRES_NAME", "server_administration_service") +
		" sslmode=disable"
	db := postgres.ConnectDB(dsn)

	// Migrate the database
	postgres.Migrate(db)

	// Initialize Redis client
	redisAddress := env.GetEnv("SERVER_REDIS_HOST", "localhost") + 
				":" + env.GetEnv("SERVER_REDIS_PORT", "6379")
	redis := redis.NewRedisClient(redisAddress)

	// Initialize Kafka Consumer Group
	brokers := []string{"localhost:9092"}
	groupID := "server_administration_group"
	topics := []string{"healthcheck_topic"}

	consumerGroup, err := kafka.NewKafkaConsumerGroup(brokers, groupID, topics)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to connect to Kafka: " + err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}

	// Initialize internal services
	serverRepository := repository.NewServerRepository(db, redis)
	serverService := service.NewServerService(serverRepository)
	serverHandler := handler.NewGrpcServerHandler(serverService)

	kafkaHandler := handler.NewServerConsumerHandler(serverService)
	consumerGroup.StartConsuming(kafkaHandler)

	// Start gRPC server
	grpcPort := env.GetEnv("SERVER_GRPC_ADMINISTRATION_PORT", "50051")
	
	logging.LogMessage("server_administration_service", "Starting gRPC server on port "+grpcPort, "INFO")
	go grpc.StartGRPCServer(serverHandler, grpcPort)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs // Wait for interrupt
	logging.LogMessage("server_administration_service", "Shutting down server...", "INFO")
	consumerGroup.Stop()
	redis.Close()
}