package main

import (
	"encoding/json"
	"healthcheck_service/infrastructure/grpc"
	"healthcheck_service/infrastructure/healthcheck"
	grpcclient "healthcheck_service/internal/grpc_client"
	"healthcheck_service/pb"
	"os"
	"path/filepath"
	"time"

	"github.com/flashhhhh/pkg/env"
	"github.com/flashhhhh/pkg/kafka"
	"github.com/flashhhhh/pkg/logging"
)

func main() {
	// Initialize logging
	currentPath, _ := os.Getwd()
	serverServiceLogPath := filepath.Join(currentPath, "logs", "healthcheck_service.log")
	logging.InitLogger("healthcheck_service", serverServiceLogPath, 10, 5, 30)

	// Initialize Kafka Producer
	kafkaProducer, err := kafka.NewKafkaProducer([]string{env.GetEnv("KAFKA_ADDRESS", "localhost:9092")})
	if err != nil {
		panic(err)
	}
	defer kafkaProducer.Close()

	topic := "healthcheck_topic"

	grpcClient, err := grpc.StartGRPCClient()
	if err != nil {
		panic(err)
	}

	client := grpcclient.NewHealthCheckClient(grpcClient)

	for {
		addressesResponse, err := client.GetAllAddresses()
		if err != nil {
			panic(err)
		}							
		for _, address := range addressesResponse.Addresses {
			go func (address *pb.AddressInfo) {
				serverID := address.ServerId
				serverAddress := address.Address
				
				// Check if the server is On or Off by pinging the address
				logging.LogMessage("healthcheck_service", "Pinging server "+serverID+" at address "+serverAddress, "INFO")
				status := healthcheck.IsHostUp(serverAddress)

				statusText := "OFF"
				if status {
					statusText = "ON"
				}
				logging.LogMessage("healthcheck_service", "Server "+serverID+" is "+statusText, "INFO")
				
				// Send the health check result to Kafka
				data := map[string]interface{}{
					"server_id":   serverID,
					"ipv4": serverAddress,
					"status":      status,
				}

				message, _ := json.Marshal(data)
				err = kafkaProducer.SendMessage(topic, message)

				if err != nil {
					panic(err)
				}

				logging.LogMessage("healthcheck_service", "Sent health check result of server " + serverID + " to Kafka topic "+topic, "INFO")
			}(address)
		}

		logging.LogMessage("healthcheck_service", "Sleep for 60 seconds...", "INFO")
		time.Sleep(60 * time.Second)
	}
}