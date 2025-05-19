package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"server_administration_service/infrastructure/elasticsearch"
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
	environment := env.GetEnv("RUNNING_ENVIRONMENT", "local")
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

	if environment == "local" {
		logging.LogMessage("server_administration_service", "Running database migrations in local environment", "INFO")
		postgres.Migrate(db)
	} else {
		logging.LogMessage("server_administration_service", "Skipping database migrations in non-local environment", "INFO")
	}

	// Initialize Redis client
	redisAddress := env.GetEnv("SERVER_REDIS_HOST", "localhost") +
			 ":" + env.GetEnv("SERVER_REDIS_PORT", "6379")
	redis := redis.NewRedisClient(redisAddress)

	elasticSearchAddress := env.GetEnv("SERVER_ELASTICSEARCH_HOST", "localhost") +
			 ":" + env.GetEnv("SERVER_ELASTICSEARCH_PORT", "9200")
	es := elasticsearch.ConnectES(elasticSearchAddress)
	
	// Initialize Kafka Consumer Group
	brokers := []string{env.GetEnv("KAFKA_HOST", "localhost") + ":" + env.GetEnv("KAFKA_PORT", "9092")}
	groupID := "server_administration_group"
	topics := []string{"healthcheck_topic"}

	logging.LogMessage("server_administration_service", "Connecting to Kafka brokers: "+brokers[0], "INFO")

	consumerGroup, err := kafka.NewKafkaConsumerGroup(brokers, groupID, topics)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to connect to Kafka: " + err.Error(), "FATAL")
		logging.LogMessage("server_administration_service", "Exiting the program...", "FATAL")
		os.Exit(1)
	}

	// Initialize internal services
	serverRepository := repository.NewServerRepository(db, redis, es)
	serverService := service.NewServerService(serverRepository)

	// Start Kafka consumer
	kafkaHandler := handler.NewServerConsumerHandler(serverService)
	consumerGroup.StartConsuming(kafkaHandler)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs // Wait for interrupt
	logging.LogMessage("server_administration_service", "Shutting down server...", "INFO")
	consumerGroup.Stop()
	redis.Close()
}