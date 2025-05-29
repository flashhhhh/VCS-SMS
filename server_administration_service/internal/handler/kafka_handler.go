package handler

import (
	"encoding/json"
	"server_administration_service/internal/service"

	"github.com/IBM/sarama"
	"github.com/flashhhhh/pkg/logging"
)

type ServerConsumerHandler struct {
	serverService service.ServerService
}

func NewServerConsumerHandler(serverService service.ServerService) *ServerConsumerHandler {
	return &ServerConsumerHandler{
		serverService: serverService,
	}
}

func (h ServerConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h ServerConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h ServerConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Create a semaphore channel to limit concurrent goroutines
	const maxWorkers = 10
	semaphore := make(chan struct{}, maxWorkers)
	
	for message := range claim.Messages() {
		// Acquire semaphore
		semaphore <- struct{}{}
		
		go func(message *sarama.ConsumerMessage) {
			// Release semaphore when done
			defer func() { <-semaphore }()
			
			logging.LogMessage("server_administration_service", "Received message: " + string(message.Value), "INFO")
			session.MarkMessage(message, "")

			// Parse the message
			var serverMessage struct {
				IPv4     string `json:"ipv4"`
				ID int `json:"id"`
				Status   bool   `json:"status"`
			}
			
			if err := json.Unmarshal(message.Value, &serverMessage); err != nil {
				logging.LogMessage("server_administration_service", "Error parsing message: "+err.Error(), "ERROR")
				return
			}
			
			// Now you can use the parsed message
			logging.LogMessage("server_administration_service", "Updating server status: "+serverMessage.IPv4, "INFO")

			status := "Off"
			if serverMessage.Status {
				status = "On"
			}
			h.serverService.UpdateServerStatus(serverMessage.ID, status)

			logging.LogMessage("server_administration_service", "Write to ES: "+serverMessage.IPv4, "INFO")
			err := h.serverService.AddServerStatus(serverMessage.ID, status)
			if err != nil {
				logging.LogMessage("server_administration_service", "Error writing to ES: "+err.Error(), "ERROR")
			}

			logging.LogMessage("server_administration_service", "Message processed: " + string(message.Value), "INFO")
		}(message)
	}

	return nil
}