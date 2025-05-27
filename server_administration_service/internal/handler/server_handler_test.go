package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/handler"
	"server_administration_service/pb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServerService is a mock implementation of service.ServerService
type MockServerService struct {
	mock.Mock
}

func (m *MockServerService) CreateServer(serverID, serverName, status, ipv4 string, port int) (int, error) {
	args := m.Called(serverID, serverName, status, ipv4, port)
	return args.Int(0), args.Error(1)
}

func (m *MockServerService) ViewServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn, order string) ([]domain.Server, error) {
	args := m.Called(serverFilter, from, to, sortedColumn, order)
	return args.Get(0).([]domain.Server), args.Error(1)
}

func (m *MockServerService) UpdateServer(serverID string, updatedData map[string]interface{}) error {
	args := m.Called(serverID, updatedData)
	return args.Error(0)
}

func (m *MockServerService) DeleteServer(serverID string) error {
	args := m.Called(serverID)
	return args.Error(0)
}

func (m *MockServerService) ImportServers(data []byte) ([]domain.Server, []domain.Server, error) {
	args := m.Called(data)
	return args.Get(0).([]domain.Server), args.Get(1).([]domain.Server), args.Error(2)
}

func (m *MockServerService) ExportServers(serverFilter *dto.ServerFilter, from, to int, sortedColumn, order string) ([]byte, error) {
	args := m.Called(serverFilter, from, to, sortedColumn, order)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockServerService) UpdateServerStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockServerService) GetAllAddresses() ([]dto.ServerAddress, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ServerAddress), args.Error(1)
}

func (m *MockServerService) AddServerStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockServerService) GetNumOnServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockServerService) GetNumServers() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockServerService) GetServerUptimeRatio(startTime, endTime time.Time) (float64, error) {
	args := m.Called(startTime, endTime)
	return args.Get(0).(float64), args.Error(1)
}

func TestCreateServer_Success(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	mockService.On("CreateServer", "server123", "Test Server", "On", "192.168.1.1", 8080).
		Return(201, nil)

	body := map[string]interface{}{
		"server_id":   "server123",
		"server_name": "Test Server",
		"status":     "On",
		"ipv4":       "192.168.1.1",
		"port":       8080,
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/create", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.CreateServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode, "Expected status code 201")
	var responseBody []byte
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Server created successfully", string(responseBody))
	
	mockService.AssertExpectations(t)
}

func TestCreateServer_Failure(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	mockService.On("CreateServer", "server123", "Test Server", "On", "192.168.1.1", 8080).
		Return(0, assert.AnError)

	body := map[string]interface{}{
		"server_id":   "server123",
		"server_name": "Test Server",
		"status":     "On",
		"ipv4":       "192.168.1.1",
		"port":       8080,
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/create", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.CreateServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "Expected status code 500")
	var responseBody []byte
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Failed to create server\n", string(responseBody), "Expected error message 'Failed to create server'")
	mockService.AssertExpectations(t)
}

func TestViewServers_Success(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	servers := []domain.Server{
		{ID: 1, ServerID: "server123", ServerName: "Test Server", Status: "On", IPv4: "192.168.1.1", Port: 8080},
		{ID: 2, ServerID: "server456", ServerName: "Another Server", Status: "Off", IPv4: "192.168.1.2", Port: 9090},
	}
	
	expectedFilter := &dto.ServerFilter{
		ServerID:   "server123",
		ServerName: "Test",
		Status:     "On",
		IPv4:       "192.168.1.1",
		Port:       8080,
	}
	
	mockService.On("ViewServers", expectedFilter, 0, 10, "server_name", "asc").Return(servers, nil)

	req := httptest.NewRequest("GET", "/servers?from=0&to=10&sort_column=server_name&sort_order=asc&server_id=server123&server_name=Test&status=On&ipv4=192.168.1.1&port=8080", nil)
	rec := httptest.NewRecorder()

	handler.ViewServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	
	var responseServers []domain.Server
	err = json.Unmarshal(responseBody, &responseServers)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}
	
	assert.Equal(t, servers, responseServers)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	
	mockService.AssertExpectations(t)
}

func TestViewServers_InvalidFrom(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("GET", "/servers?from=invalid&to=10", nil)
	rec := httptest.NewRecorder()

	handler.ViewServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Invalid 'from' query parameter\n", string(responseBody))
}

func TestViewServers_InvalidTo(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("GET", "/servers?from=0&to=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ViewServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Invalid 'to' query parameter\n", string(responseBody))
}

func TestViewServers_InvalidPort(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("GET", "/servers?from=0&to=10&port=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ViewServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Invalid port\n", string(responseBody))
}

func TestViewServers_ServiceError(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	expectedFilter := &dto.ServerFilter{Port: -1}
	mockService.On("ViewServers", expectedFilter, 0, 10, "", "").Return([]domain.Server{}, assert.AnError)

	req := httptest.NewRequest("GET", "/servers?from=0&to=10", nil)
	rec := httptest.NewRecorder()

	handler.ViewServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Failed to view servers\n", string(responseBody))
	
	mockService.AssertExpectations(t)
}

func TestUpdateServer_Success(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	updatedData := map[string]interface{}{
		"server_name": "Updated Server",
		"status":      "Off",
		"ipv4":        "192.168.1.2",
		"port":        9090,
	}
	
	mockService.On("UpdateServer", "server123", updatedData).Return(nil)

	body := map[string]interface{}{
		"server_name": "Updated Server",
		"status":      "Off",
		"ipv4":        "192.168.1.2",
		"port":        9090.0,
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/update?server_id=server123", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.UpdateServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode, "Expected status code 200")
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Server updated successfully", string(responseBody))
	
	mockService.AssertExpectations(t)
}

func TestUpdateServer_MissingID(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	body := map[string]interface{}{
		"server_name": "Updated Server",
		"status":      "Off",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/update", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.UpdateServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Expected status code 400")
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Server ID is required\n", string(responseBody))
}

func TestUpdateServer_InvalidBody(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("PUT", "/update?server_id=server123", bytes.NewBufferString("invalid json"))
	rec := httptest.NewRecorder()

	handler.UpdateServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Expected status code 400")
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Invalid request body\n", string(responseBody))
}

func TestUpdateServer_ServiceError(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	updatedData := map[string]interface{}{
		"server_name": "Updated Server",
	}
	
	mockService.On("UpdateServer", "server123", updatedData).Return(assert.AnError)

	body := map[string]interface{}{
		"server_name": "Updated Server",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/update?server_id=server123", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.UpdateServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "Expected status code 500")
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Failed to update server\n", string(responseBody))
	
	mockService.AssertExpectations(t)
}

func TestDeleteServer_Success(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	serverID := "server123"
	mockService.On("DeleteServer", serverID).Return(nil)

	req := httptest.NewRequest("DELETE", "/delete?server_id="+serverID, nil)
	rec := httptest.NewRecorder()

	handler.DeleteServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode, "Expected status code 200")
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Server deleted successfully", string(responseBody))
	
	mockService.AssertExpectations(t)
}

func TestDeleteServer_MissingID(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("DELETE", "/delete", nil)
	rec := httptest.NewRecorder()

	handler.DeleteServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Expected status code 400")
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Server ID is required\n", string(responseBody))
}

func TestDeleteServer_ServiceError(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	serverID := "nonexistent"
	mockService.On("DeleteServer", serverID).Return(assert.AnError)

	req := httptest.NewRequest("DELETE", "/delete?server_id="+serverID, nil)
	rec := httptest.NewRecorder()

	handler.DeleteServer(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusNotFound, res.StatusCode, "Expected status code 404")
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Invalid server ID\n", string(responseBody))
	
	mockService.AssertExpectations(t)
}

func TestExportServers_Success(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	excelBytes := []byte("mock excel data")
	expectedFilter := &dto.ServerFilter{
		ServerID:   "server123",
		ServerName: "TestServer",
		Status:     "On",
		IPv4:       "192.168.1.1",
		Port:       8080,
	}
	
	mockService.On("ExportServers", expectedFilter, 0, 10, "server_name", "asc").Return(excelBytes, nil)

	req := httptest.NewRequest("GET", "/export?from=0&to=10&sort_column=server_name&sort_order=asc&server_id=server123&server_name=TestServer&status=On&ipv4=192.168.1.1&port=8080", nil)
	rec := httptest.NewRecorder()

	handler.ExportServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", res.Header.Get("Content-Type"))
	assert.Contains(t, res.Header.Get("Content-Disposition"), "attachment; filename=servers_")
	assert.Contains(t, res.Header.Get("Content-Disposition"), ".xlsx")
	assert.Contains(t, res.Header.Get("File-Name"), "servers_")
	assert.Contains(t, res.Header.Get("File-Name"), ".xlsx")
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, excelBytes, responseBody)
	
	mockService.AssertExpectations(t)
}

func TestExportServers_InvalidFrom(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("GET", "/export?from=invalid&to=10", nil)
	rec := httptest.NewRecorder()

	handler.ExportServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Invalid 'from' query parameter\n", string(responseBody))
}

func TestExportServers_InvalidTo(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("GET", "/export?from=0&to=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ExportServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Invalid 'to' query parameter\n", string(responseBody))
}

func TestExportServers_InvalidPort(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("GET", "/export?from=0&to=10&port=invalid", nil)
	rec := httptest.NewRecorder()

	handler.ExportServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Invalid port\n", string(responseBody))
}

func TestExportServers_ServiceError(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	expectedFilter := &dto.ServerFilter{Port: -1}
	mockService.On("ExportServers", expectedFilter, 0, 10, "", "").Return([]byte{}, assert.AnError)

	req := httptest.NewRequest("GET", "/export?from=0&to=10", nil)
	rec := httptest.NewRecorder()

	handler.ExportServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Failed to export servers\n", string(responseBody))
	
	mockService.AssertExpectations(t)
}

func TestGetAllAddresses_Success(t *testing.T) {
	mockService := new(MockServerService)
	grpcHandler := handler.NewGrpcServerHandler(mockService)

	addresses := []dto.ServerAddress{
		{ID: 1, IPv4: "192.168.1.1", Port: 8080},
		{ID: 2, IPv4: "192.168.1.2", Port: 9090},
	}
	mockService.On("GetAllAddresses").Return(addresses, nil)

	req := &pb.EmptyRequest{}
	response, err := grpcHandler.GetAllAddresses(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Addresses))
	
	mockService.AssertExpectations(t)
}

func TestGetAllAddresses_ServiceError(t *testing.T) {
	mockService := new(MockServerService)
	grpcHandler := handler.NewGrpcServerHandler(mockService)

	mockService.On("GetAllAddresses").Return(nil, assert.AnError)

	req := &pb.EmptyRequest{}
	response, err := grpcHandler.GetAllAddresses(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, response)
	
	mockService.AssertExpectations(t)
}

func TestGetServerInformation_Success(t *testing.T) {
	mockService := new(MockServerService)
	grpcHandler := handler.NewGrpcServerHandler(mockService)
	
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	
	mockService.On("GetNumOnServers").Return(3, nil)
	mockService.On("GetNumServers").Return(5, nil)
	mockService.On("GetServerUptimeRatio", 
		mock.MatchedBy(func(st time.Time) bool { 
			return st.Unix() == startTime.Unix() 
		}),
		mock.MatchedBy(func(et time.Time) bool { 
			return et.Unix() == endTime.Unix() 
		})).Return(0.75, nil)
	
	req := &pb.GetServerInformationRequest{
		StartTime: startTime.Unix(),
		EndTime: endTime.Unix(),
	}
	
	response, err := grpcHandler.GetServerInformation(context.Background(), req)
	
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(5), response.NumServers)
	assert.Equal(t, int64(3), response.NumOnServers)
	assert.Equal(t, int64(2), response.NumOffServers)
	assert.Equal(t, float32(0.75), response.MeanUptimeRatio)
	
	mockService.AssertExpectations(t)
}

func TestGetServerInformation_UptimeRatioError(t *testing.T) {
	mockService := new(MockServerService)
	grpcHandler := handler.NewGrpcServerHandler(mockService)
	
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	
	mockService.On("GetNumOnServers").Return(3, nil)
	mockService.On("GetNumServers").Return(5, nil)
	mockService.On("GetServerUptimeRatio", 
		mock.MatchedBy(func(st time.Time) bool { 
			return st.Unix() == startTime.Unix() 
		}),
		mock.MatchedBy(func(et time.Time) bool { 
			return et.Unix() == endTime.Unix() 
		})).Return(0.0, assert.AnError)
	
	req := &pb.GetServerInformationRequest{
		StartTime: startTime.Unix(),
		EndTime: endTime.Unix(),
	}
	
	response, err := grpcHandler.GetServerInformation(context.Background(), req)
	
	assert.Error(t, err)
	assert.Nil(t, response)
	
	mockService.AssertExpectations(t)
}

func TestGetServerInformation_EmptyTimestamps(t *testing.T) {
	mockService := new(MockServerService)
	grpcHandler := handler.NewGrpcServerHandler(mockService)
	
	mockService.On("GetNumOnServers").Return(3, nil)
	mockService.On("GetNumServers").Return(5, nil)
	mockService.On("GetServerUptimeRatio", 
		mock.MatchedBy(func(st time.Time) bool { 
			return st.Unix() == 0 
		}),
		mock.MatchedBy(func(et time.Time) bool { 
			return et.Unix() == 0 
		})).Return(0.8, nil)
	
	req := &pb.GetServerInformationRequest{
		StartTime: 0,
		EndTime: 0,
	}
	
	response, err := grpcHandler.GetServerInformation(context.Background(), req)
	
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(5), response.NumServers)
	assert.Equal(t, int64(3), response.NumOnServers)
	assert.Equal(t, int64(2), response.NumOffServers)
	assert.Equal(t, float32(0.8), response.MeanUptimeRatio)
	
	mockService.AssertExpectations(t)
}

func TestImportServers_Success(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	importedServers := []domain.Server{
		{ID: 1, ServerID: "server123", ServerName: "Test Server", Status: "On", IPv4: "192.168.1.1", Port: 8080},
	}
	nonImportedServers := []domain.Server{
		{ID: 0, ServerID: "server456", ServerName: "Invalid Server", Status: "Unknown", IPv4: "invalid", Port: 0},
	}
	
	// Mock file content
	fileContent := []byte("mock excel data")
	
	mockService.On("ImportServers", fileContent).Return(importedServers, nonImportedServers, nil)

	// Create multipart form with file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("servers_file", "test.xlsx")
	part.Write(fileContent)
	writer.Close()

	req := httptest.NewRequest("POST", "/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.ImportServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	
	var response map[string]interface{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}
	
	assert.Contains(t, response, "imported_servers")
	assert.Contains(t, response, "non_imported_servers")
	
	mockService.AssertExpectations(t)
}

func TestImportServers_NoFile(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)

	req := httptest.NewRequest("POST", "/import", nil)
	rec := httptest.NewRecorder()

	handler.ImportServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Failed to get file from request\n", string(responseBody))
}

func TestImportServers_ServiceError(t *testing.T) {
	mockService := new(MockServerService)
	handler := handler.NewServerHandler(mockService)
	
	// Mock file content
	fileContent := []byte("mock excel data")
	
	mockService.On("ImportServers", fileContent).Return([]domain.Server{}, []domain.Server{}, assert.AnError)

	// Create multipart form with file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("servers_file", "test.xlsx")
	part.Write(fileContent)
	writer.Close()

	req := httptest.NewRequest("POST", "/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.ImportServers(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	assert.Equal(t, "Failed to import servers\n", string(responseBody))
	
	mockService.AssertExpectations(t)
}