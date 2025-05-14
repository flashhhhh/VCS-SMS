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