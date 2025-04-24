package handler

import (
	"encoding/json"
	"net/http"
	"server_administration_service/internal/service"

	"github.com/flashhhhh/pkg/logging"
)

type ServerHandler interface {
	CreateServer(w http.ResponseWriter, r *http.Request)
	UpdateServer(w http.ResponseWriter, r *http.Request)
	DeleteServer(w http.ResponseWriter, r *http.Request)
}

type serverHandler struct {
	service service.ServerService
}

func NewServerHandler(service service.ServerService) ServerHandler {
	return &serverHandler{
		service: service,
	}
}

func (h *serverHandler) CreateServer(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to decode request body for request CreateServer: "+err.Error(), "ERROR")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	serverID, _ := requestBody["server_id"].(string)
	serverName, _ := requestBody["server_name"].(string)
	status, _ := requestBody["status"].(string)
	ipAddress, _ := requestBody["ipv4"].(string)
	portFloat, ok := requestBody["port"].(float64)

	port := 80 // default port
	if ok {
		port = int(portFloat)
	}
	
	err = h.service.CreateServer(serverID, serverName, status, ipAddress, port)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to create server: "+err.Error(), "ERROR")
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
		return
	}

	logging.LogMessage("server_administration_service", "Server created successfully with ID: "+serverID, "INFO")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Server created successfully"))
}

func (h *serverHandler) UpdateServer(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	if serverID == "" {
		logging.LogMessage("server_administration_service", "Server ID is required for update", "ERROR")
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to decode request body for request UpdateServer: "+err.Error(), "ERROR")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedData := make(map[string]interface{})
	serverName, existed := requestBody["server_name"].(string)
	if existed {
		updatedData["server_name"] = serverName
	}

	status, existed := requestBody["status"].(string)
	if existed {
		updatedData["status"] = status
	}

	ipAddress, existed := requestBody["ipv4"].(string)
	if existed {
		updatedData["ipv4"] = ipAddress
	}

	portFloat, existed := requestBody["port"].(float64)
	if existed {
		updatedData["port"] = int(portFloat)
	}

	err = h.service.UpdateServer(serverID, updatedData)
	if err != nil {
		logging.LogMessage("server_administration_service", "Failed to update server: "+err.Error(), "ERROR")
		http.Error(w, "Failed to update server", http.StatusInternalServerError)
		return
	}

	logging.LogMessage("server_administration_service", "Server updated successfully with ID: "+serverID, "INFO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server updated successfully"))
}

func (h *serverHandler) DeleteServer(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	if serverID == "" {
		logging.LogMessage("server_administration_service", "Server ID is required for deletion", "ERROR")
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteServer(serverID)
	if err != nil {
		logging.LogMessage("server_administration_service", "Invalid server ID: "+serverID+" - "+err.Error(), "ERROR")
		http.Error(w, "Invalid server ID", http.StatusNotFound)
		return
	}

	logging.LogMessage("server_administration_service", "Server deleted successfully with ID: "+serverID, "INFO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server deleted successfully"))
}