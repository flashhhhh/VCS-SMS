// user_service_test.go
package service_test

import (
	"context"
	"errors"
	"testing"

	// "time"
	"user_service/internal/domain"
	"user_service/internal/service"

	// "github.com/flashhhhh/pkg/hash"
	// "github.com/flashhhhh/pkg/jwt"
	"github.com/flashhhhh/pkg/hash"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *mockUserRepo) Login(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func TestLogin_Successs(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	username := "testuser"
	password := "testpassword"
	hashedPassword := hash.HashString(password)

	user := &domain.User{
		ID:       "1",
		Username: username,
		Password: hashedPassword,
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	mockRepo.On("Login", mock.Anything, username).Return(user, nil).Once()

	token, err := userService.Login(context.Background(), username, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected token, got empty string")
	}
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	username := "testuser"
	password := "wrongpassword"
	hashedPassword := hash.HashString("testpassword")

	user := &domain.User{
		ID:       "1",
		Username: username,
		Password: hashedPassword,
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	mockRepo.On("Login", mock.Anything, username).Return(user, nil).Once()
	_, err := userService.Login(context.Background(), username, password)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "Invalid password" {
		t.Fatalf("expected 'Invalid password' error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestLogin_NotFound(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	username := "testuser"
	password := "testpassword"

	mockRepo.On("Login", mock.Anything, username).Return(nil, errors.New("user not found")).Once()

	_, err := userService.Login(context.Background(), username, password)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "user not found" {
		t.Fatalf("expected 'user not found' error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}