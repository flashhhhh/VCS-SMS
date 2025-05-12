// user_service_test.go
package service_test

import (
	"context"
	"errors"
	"testing"

	"user_service/internal/domain"
	"user_service/internal/service"

	"github.com/flashhhhh/pkg/hash"
	"github.com/google/uuid"
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	user := &domain.User{
		Username: "testuser",
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	expectedID := uuid.New().String()
	
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Run(func(args mock.Arguments) {
		capturedUser := args.Get(1).(*domain.User)
		capturedUser.ID = expectedID
	}).Return(expectedID, nil).Once()

	userID, err := userService.CreateUser(context.Background(), user.Username, "testpassword", user.Name, user.Email, user.Role)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if userID != expectedID {
		t.Fatalf("expected user ID %s, got %s", expectedID, userID)
	}

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Failure(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	user := &domain.User{
		Username: "testuser",
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return("", errors.New("failed to create user")).Once()
	_, err := userService.CreateUser(context.Background(), user.Username, "testpassword", user.Name, user.Email, user.Role)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "failed to create user" {
		t.Fatalf("expected 'failed to create user' error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
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

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	userID := "1"
	user := &domain.User{
		ID:       userID,
		Username: "testuser",
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil).Once()
	result, err := userService.GetUserByID(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatalf("expected user, got nil")
	}
	if result.ID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, result.ID)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Failure(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	userID := "1"

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(nil, errors.New("user not found")).Once()
	_, err := userService.GetUserByID(context.Background(), userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "User 1 not found: user not found" {
		t.Fatalf("expected 'User 1 not found: user not found' error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	users := []*domain.User{
		{
			ID:	   "1",
			Username: "testuser1",
			Name:     "Test User 1",
			Email:    "testuser1@gmail.com",
			Role:     "user",
		},
		{
			ID:	   "2",
			Username: "testuser2",
			Name:     "Test User 2",
			Email:    "testuser2@gmail.com",
			Role:     "admin",
		},
	}
	mockRepo.On("GetAllUsers", mock.Anything).Return(users, nil).Once()

	result, err := userService.GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatalf("expected users, got nil")
	}
	if len(result) != len(users) {
		t.Fatalf("expected %d users, got %d", len(users), len(result))
	}
	for i, user := range users {
		if result[i].ID != user.ID {
			t.Fatalf("expected user ID %s, got %s", user.ID, result[i].ID)
		}
	}
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Failure(t *testing.T) {
	mockRepo := new(mockUserRepo)
	userService := service.NewUserService(mockRepo)

	mockRepo.On("GetAllUsers", mock.Anything).Return(nil, errors.New("failed to get users")).Once()
	_, err := userService.GetAllUsers(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "Failed to retrieve all users: failed to get users" {
		t.Fatalf("expected 'Failed to retrieve all users: failed to get users' error, got %v", err)
	}
	mockRepo.AssertExpectations(t)
}