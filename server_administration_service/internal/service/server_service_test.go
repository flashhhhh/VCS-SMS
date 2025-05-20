package service_test

import (
	"errors"
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type mockServerRepo struct {
	mock.Mock
}

func (m *mockServerRepo) CreateServer(server *domain.Server) (int, error) {
	args := m.Called(server)
	return args.Int(0), args.Error(1)
}

func (m *mockServerRepo) CreateServers(servers []domain.Server) ([]domain.Server, []domain.Server, error) {
	args := m.Called(servers)
	return args.Get(0).([]domain.Server), args.Get(1).([]domain.Server), args.Error(2)
}

func (m *mockServerRepo) ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn string, order string) ([]domain.Server, error) {
	args := m.Called(serverFilter, from, to, sortedColumn, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Server), args.Error(1)
}

func (m *mockServerRepo) UpdateServer(serverID string, updatedData map[string]interface{}) error {
	args := m.Called(serverID, updatedData)
	return args.Error(0)
}

func (m *mockServerRepo) DeleteServer(serverID string) error {
	args := m.Called(serverID)
	return args.Error(0)
}

func (m *mockServerRepo) UpdateServerStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *mockServerRepo) GetAllAddresses() ([]dto.ServerAddress, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ServerAddress), args.Error(1)
}

func (m *mockServerRepo) AddServerStatus(serverID int, status string) error {
	args := m.Called(serverID, status)
	return args.Error(0)
}

func (m *mockServerRepo) GetNumOnServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *mockServerRepo) GetNumServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}
func (m *mockServerRepo) GetServerUptimeRatio(startTime, endTime time.Time) (float64, error) {
	args := m.Called(startTime, endTime)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockServerRepo) SyncServerStatus() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateServer_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	server := &domain.Server{
		ServerID:   "server123",
		ServerName: "Test Server",
		Status:     "On",
		IPv4:       "192.168.1.1",
		Port:       8080,
	}
	mockRepo.On("CreateServer", server).Return(1, nil)

	id, err := serverService.CreateServer("server123", "Test Server", "On", "192.168.1.1", 8080)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if id != 1 {
		t.Errorf("Expected server ID 1, got %d", id)
	}
	mockRepo.AssertExpectations(t)
}

func TestCreateServer_Failure(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	server := &domain.Server{
		ServerID:   "server123",
		ServerName: "Test Server",
		Status:     "On",
		IPv4:       "192.168.1.1",
		Port:       8080,
	}
	mockRepo.On("CreateServer", server).Return(0, errors.New("server creation failed"))

	id, err := serverService.CreateServer("server123", "Test Server", "On", "192.168.1.1", 8080)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if id != 0 {
		t.Errorf("Expected server ID 0, got %d", id)
	}
	mockRepo.AssertExpectations(t)
}

func TestViewServers_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	serverFilter := &dto.ServerFilter{
		Status: "On",
	}
	servers := []domain.Server{
		{
			ServerID:   "server123",
			ServerName: "Test Server",
			Status:     "On",
			IPv4:       "192.168.1.1",
			Port:       8080,
		},
		{
			ServerID:   "server124",
			ServerName: "Test Server 2",
			Status:     "On",
			IPv4:       "192.168.1.2",
			Port:       8081,
		},
	}
	mockRepo.On("ViewServers", serverFilter, 0, 10, "ServerID", "asc").Return(servers, nil)

	result, err := serverService.ViewServers(serverFilter, 0, 10, "ServerID", "asc")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(result))
	}
	if result[0].ServerID != "server123" {
		t.Errorf("Expected server ID 'server123', got %s", result[0].ServerID)
	}
	if result[1].ServerID != "server124" {
		t.Errorf("Expected server ID 'server124', got %s", result[1].ServerID)
	}
	mockRepo.AssertExpectations(t)
}

func TestViewServers_Failure(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	serverFilter := &dto.ServerFilter{
		Status: "On",
	}
	mockRepo.On("ViewServers", serverFilter, 0, 10, "ServerID", "asc").Return(nil, errors.New("failed to fetch servers"))

	result, err := serverService.ViewServers(serverFilter, 0, 10, "ServerID", "asc")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	mockRepo.AssertExpectations(t)
}

func TestUpdateServer_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	updatedData := map[string]interface{}{
		"Status": "Off",
	}
	mockRepo.On("UpdateServer", "server123", updatedData).Return(nil)

	err := serverService.UpdateServer("server123", updatedData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestDeleteServer_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	mockRepo.On("DeleteServer", "server123").Return(nil)
	err := serverService.DeleteServer("server123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestUpdateServerStatus_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	mockRepo.On("UpdateServerStatus", 1, "On").Return(nil)
	err := serverService.UpdateServerStatus(1, "On")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetAllAddresses_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	addresses := []dto.ServerAddress{
		{
			ID: 1,
			IPv4:     "192.168.1.1",
			Port:     8080,
		},
		{
			ID: 2,
			IPv4:     "192.168.1.2",
			Port:     8081,
		},
	}
	mockRepo.On("GetAllAddresses").Return(addresses, nil)

	result, err := serverService.GetAllAddresses()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 addresses, got %d", len(result))
	}

	mockRepo.AssertExpectations(t)
}

func TestGetAllAddresses_Failure(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	mockRepo.On("GetAllAddresses").Return(nil, errors.New("failed to fetch addresses"))

	result, err := serverService.GetAllAddresses()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetNumOnServers_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	mockRepo.On("GetNumOnServers").Return(5, nil)

	result, err := serverService.GetNumOnServers()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 5 {
		t.Errorf("Expected 5 servers, got %d", result)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetNumServers(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	mockRepo.On("GetNumServers").Return(10, nil)

	result, err := serverService.GetNumServers()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 10 {
		t.Errorf("Expected 10 servers, got %d", result)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetServerUptimeRatio_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()
	mockRepo.On("GetServerUptimeRatio", startTime, endTime).Return(0.95, nil)

	result, err := serverService.GetServerUptimeRatio(startTime, endTime)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 0.95 {
		t.Errorf("Expected uptime ratio 0.95, got %f", result)
	}
	mockRepo.AssertExpectations(t)
}

func TestAddServerStatus(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	mockRepo.On("AddServerStatus", 1, "On").Return(nil)

	err := serverService.AddServerStatus(1, "On")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestExportServers_Success(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	servers := []domain.Server{
		{
			ServerID:   "server123",
			ServerName: "Test Server",
			Status:     "On",
			IPv4:       "192.168.1.1",
			Port:       8080,
		},
		{
			ServerID:   "server124",
			ServerName: "Test Server 2",
			Status:     "On",
			IPv4:       "192.168.1.2",
			Port:       8081,
		},
	}
	
	filter := &dto.ServerFilter{
		ServerID:   "server123",
		ServerName: "Test Server",
		Status:     "On",
		IPv4:       "192.168.1.1",
		Port:       8080,
	}
	mockRepo.On("ViewServers", filter, 0, 10, "ServerID", "asc").Return(servers, nil)
	
	_, err := serverService.ExportServers(filter, 0, 10, "ServerID", "asc")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestExportServers_FailRetrieval(t *testing.T) {
	mockRepo := new(mockServerRepo)
	serverService := service.NewServerService(mockRepo)
	
	filter := &dto.ServerFilter{
		ServerID:   "server123",
		ServerName: "Test Server",
		Status:     "On",
		IPv4:       "192.168.1.1",
		Port:       8080,
	}
	mockRepo.On("ViewServers", filter, 0, 10, "ServerID", "asc").Return(nil, errors.New("failed to fetch servers"))

	_, err := serverService.ExportServers(filter, 0, 10, "ServerID", "asc")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	mockRepo.AssertExpectations(t)
}