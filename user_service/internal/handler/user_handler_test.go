package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"user_service/internal/domain"
	"user_service/internal/handler"

	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) CreateUser(ctx context.Context, username, password, name, email, role string) (string, error) {
	args := m.Called(ctx, username, password, name, email, role)
	return args.String(0), args.Error(1)
}

func (m *mockUserService) Login(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

func (m *mockUserService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func TestCreateUserHandler_Success(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	mockService.On("CreateUser", mock.Anything, "testuser", "testpass", "Test User", "testuser@gmail.com", "admin").Return("user123", nil)
	body := map[string]string{
		"username": "testuser",
		"password": "testpass",
		"name":     "Test User",
		"email":    "testuser@gmail.com",
		"role":     "admin",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/create", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()
	handler.CreateUser(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}
	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	if response["message"] != "User created successfully" {
		t.Errorf("Expected message 'User created successfully', got %s", response["message"])
	}
	if response["userID"] != "user123" {
		t.Errorf("Expected userID 'user123', got %s", response["userID"])
	}
	mockService.AssertExpectations(t)
}

func TestCreateUserHandler_Failure(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	mockService.On("CreateUser", mock.Anything, "testuser", "testpass", "Test User", "testuser@gmail.com", "admin").Return("", errors.New("Internal server error"))
	body := map[string]string{
		"username": "testuser",
		"password": "testpass",
		"name":     "Test User",
		"email":    "testuser@gmail.com",
		"role":     "admin",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/create", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()
	handler.CreateUser(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", res.StatusCode)
	}

	mockService.AssertExpectations(t)
}

func TestLoginHandler_Success(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	mockService.On("Login", mock.Anything, "testuser", "testpass").Return("token123", nil)

	body := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response["message"] != "User logged in successfully" {
		t.Errorf("Expected message 'Login successful', got %s", response["message"])
	}

	if response["token"] != "token123" {
		t.Errorf("Expected token 'token123', got %s", response["token"])
	}

	mockService.AssertExpectations(t)
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	mockService.On("Login", mock.Anything, "testuser", "wrongpass").Return("", errors.New("Invalid password"))

	body := map[string]string{
		"username": "testuser",
		"password": "wrongpass",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	
	if res.StatusCode != 401 {
		t.Errorf("Expected status code 401, got %d", res.StatusCode)
	}

	mockService.AssertExpectations(t)
}

func TestLoginHandler_Failure(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	mockService.On("Login", mock.Anything, "testuser", "testpass").Return("", errors.New("Internal server error"))

	body := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	
	if res.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", res.StatusCode)
	}

	mockService.AssertExpectations(t)
}

func TestGetUserByIDHandler_Success(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	user := &domain.User{
		ID: "user123",
		Name: "Test User",
		Username: "testuser",
		Password: "testpass",
		Email: "testuser@gmail.com",
		Role: "admin",
	}
	mockService.On("GetUserByID", mock.Anything, "user123").Return(user, nil)

	req := httptest.NewRequest("GET", "/user/?userID=user123", nil)
	rec := httptest.NewRecorder()
	handler.GetUserByID(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
	if response["user"] == nil {
		t.Errorf("Expected user object, got nil")
	}
	
	mockService.AssertExpectations(t)
}

func TestGetUserByIDHandler_WrongUserID(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	mockService.On("GetUserByID", mock.Anything, "user123").Return(nil, errors.New("User " + "user123" + " not found"))
	req := httptest.NewRequest("GET", "/user/?userID=user123", nil)
	rec := httptest.NewRecorder()
	handler.GetUserByID(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", res.StatusCode)
	}

	mockService.AssertExpectations(t)
}

func TestGetUserByIDHandler_Failure(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	mockService.On("GetUserByID", mock.Anything, "user123").Return(nil, errors.New("Internal server error"))
	req := httptest.NewRequest("GET", "/getUserByID/?userID=user123", nil)
	rec := httptest.NewRecorder()
	handler.GetUserByID(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", res.StatusCode)
	}
	mockService.AssertExpectations(t)
}

func TestGetAllUsersHandler_Success(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	users := []*domain.User{
		{
			ID: "user123",
			Name: "Test User",
			Username: "testuser",
			Password: "testpass",
			Email: "testuser@gmail.com",
			Role: "admin",
		},
		{
			ID: "user456",
			Name: "Another User",
			Username: "anotheruser",
			Password: "anotherpass",
			Email: "anotheruser@gmail.com",
			Role: "user",
		},
	}
	mockService.On("GetAllUsers", mock.Anything).Return(users, nil)

	req := httptest.NewRequest("GET", "/getAllUsers", nil)
	rec := httptest.NewRecorder()
	handler.GetAllUsers(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	
	if res.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	
	mockService.AssertExpectations(t)
}

func TestGetAllUsersHandler_Failure(t *testing.T) {
	mockService := new(mockUserService)
	handler := handler.NewUserHandler(mockService)

	mockService.On("GetAllUsers", mock.Anything).Return(nil, errors.New("Internal server error"))

	req := httptest.NewRequest("GET", "/getAllUsers", nil)
	rec := httptest.NewRecorder()
	handler.GetAllUsers(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	
	if res.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", res.StatusCode)
	}
	
	mockService.AssertExpectations(t)
}