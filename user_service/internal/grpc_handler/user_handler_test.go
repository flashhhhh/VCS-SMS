package grpchandler_test

import (
	"context"
	"errors"
	"testing"
	"user_service/internal/domain"
	grpchandler "user_service/internal/grpc_handler"
	"user_service/pb"

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

func TestCreateUser_Success(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	mockUserService.On("CreateUser", mock.Anything, "testuser", "password", "Test User", "testuser@gmail.com", "user").Return("12345", nil)

	req := &pb.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	resp, err := grpcHandler.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.UserID != "12345" {
		t.Fatalf("expected userID to be '12345', got %v", resp.UserID)
	}

	mockUserService.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	mockUserService.On("Login", mock.Anything, "testuser", "password").Return("token123", nil)

	req := &pb.LoginRequest{
		Username: "testuser",
		Password: "password",
	}

	resp, err := grpcHandler.Login(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Token != "token123" {
		t.Fatalf("expected token to be 'token123', got %v", resp.Token)
	}

	mockUserService.AssertExpectations(t)
}

func TestCreateUser_UserAlreadyExists(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	mockUserService.On("CreateUser", mock.Anything, "testuser", "password", "Test User", "testuser@gmail.com", "user").
					Return("", errors.New("user already exists"))
	
	req := &pb.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	_, err := grpcHandler.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "user already exists" {
		t.Fatalf("expected 'user already exists' error, got %v", err)
	}
	mockUserService.AssertExpectations(t)
}

func TestCreateUser_Failure(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	mockUserService.On("CreateUser", mock.Anything, "testuser", "password", "Test User", "testuser@gmail.com", "user").
					Return("", errors.New("failed to create user"))

	req := &pb.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	_, err := grpcHandler.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "failed to create user" {
		t.Fatalf("expected 'failed to create user' error, got %v", err)
	}
	mockUserService.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	mockUserService.On("Login", mock.Anything, "testuser", "wrongpassword").Return("", errors.New("invalid password"))

	req := &pb.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err := grpcHandler.Login(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	mockUserService.AssertExpectations(t)
}

func TestLogin_Failure(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	mockUserService.On("Login", mock.Anything, "testuser", "password").Return("", errors.New("login failed"))

	req := &pb.LoginRequest{
		Username: "testuser",
		Password: "password",
	}

	_, err := grpcHandler.Login(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	mockUserService.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	user := &domain.User{
		ID:    "12345",
		Name:  "Test User",
		Email: "testuser@gmail.com",
		Role:  "user",
	}
	mockUserService.On("GetUserByID", mock.Anything, "12345").Return(user, nil)
	req := &pb.IDRequest{
		Id: "12345",
	}

	_, err := grpcHandler.GetUserByID(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	mockUserService.AssertExpectations(t)
}

func TestGetUserByID_Failure(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	mockUserService.On("GetUserByID", mock.Anything, "12345").Return(nil, errors.New("user not found"))
	req := &pb.IDRequest{
		Id: "12345",
	}
	_, err := grpcHandler.GetUserByID(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "user not found" {
		t.Fatalf("expected 'user not found' error, got %v", err)
	}
	mockUserService.AssertExpectations(t)
}

func TestGetAllUsers_Success(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	users := []*domain.User{
		{ID: "12345", Name: "Test User", Email: "testuser1@gmail.com", Role: "user"},
		{ID: "67890", Name: "Another User", Email: "testuser2@gmail.com", Role: "admin"},
	}
	mockUserService.On("GetAllUsers", mock.Anything).Return(users, nil)
	req := &pb.EmptyRequest{}
	_, err := grpcHandler.GetAllUsers(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	mockUserService.AssertExpectations(t)
}

func TestGetAllUsers_Failure(t *testing.T) {
	mockUserService := new(mockUserService)
	grpcHandler := grpchandler.NewGrpcUserHandler(mockUserService)

	mockUserService.On("GetAllUsers", mock.Anything).Return(nil, errors.New("failed to fetch users"))
	req := &pb.EmptyRequest{}
	_, err := grpcHandler.GetAllUsers(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "failed to fetch users" {
		t.Fatalf("expected 'failed to fetch users' error, got %v", err)
	}
	mockUserService.AssertExpectations(t)
}