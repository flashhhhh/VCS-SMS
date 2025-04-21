// user_service_test.go
package service_test

import (
	"context"
	"errors"
	"testing"
	"user_service/internal/domain"
	"user_service/internal/service"

	"github.com/flashhhhh/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) Login(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)

	userArg := args.Get(0)
	if userArg == nil {
		return nil, args.Error(1)
	}
	return userArg.(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	
	userArg := args.Get(0)
	if userArg == nil {
		return nil, args.Error(1)
	}
	return userArg.(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	userArg := args.Get(0)
	if userArg == nil {
		return nil, args.Error(1)
	}
	return userArg.([]*domain.User), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	input := &domain.User{
		Username: "jane",
		Password: "plaintext", // Will be hashed
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Role:     "admin",
	}
	// Use a matcher instead of exact match since ID and hashed password are dynamic
	repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		return user.Username == input.Username &&
			user.Email == input.Email &&
			user.Name == input.Name &&
			user.Role == input.Role &&
			user.Password != input.Password // Should be hashed
	})).Return("some-uuid", nil)

	id, err := s.CreateUser(context.Background(), input.Username, input.Password, input.Name, input.Email, input.Role)
	assert.NoError(t, err)
	assert.Equal(t, "some-uuid", id)
}

func TestCreateUser_Failure(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	repo.On("CreateUser", mock.Anything, mock.Anything).Return("", errors.New("db error"))

	user := &domain.User{
		Username: "jane",
		Password: "plaintext",
		Name:     "Jane Doe",
		Email:    "",
		Role:     "admin",
	}

	_, err := s.CreateUser(context.Background(), user.Username, user.Password, user.Name, user.Email, user.Role)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestLogin_Success(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	hashed := hash.HashString("secret")
	user := &domain.User{
		ID:       "user123",
		Username: "john",
		Password: hashed,
		Name:     "John",
		Email:    "john@example.com",
		Role:     "user",
	}
	repo.On("Login", mock.Anything, "john").Return(user, nil)

	token, err := s.Login(context.Background(), "john", "secret")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLogin_DBError(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	repo.On("Login", mock.Anything, "john").Return(nil, errors.New("db error"))

	_, err := s.Login(context.Background(), "john", "secret")
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	hashed := hash.HashString("secret") // "secret"
	user := &domain.User{Username: "john", Password: hashed}
	repo.On("Login", mock.Anything, "john").Return(user, nil)

	_, err := s.Login(context.Background(), "john", "wrongpass")
	assert.Error(t, err)
	assert.Equal(t, "Invalid password", err.Error())
}

func TestGetUserByID_Success(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	user := &domain.User{
		ID:       "user123",
		Username: "john",
		Name:     "John",
		Email:    "",
		Role:     "user",
	}

	repo.On("GetUserByID", mock.Anything, "user123").Return(user, nil)
	result, err := s.GetUserByID(context.Background(), "user123")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	repo.On("GetUserByID", mock.Anything, "nonexistent").Return(&domain.User{}, errors.New("db error"))

	_, err := s.GetUserByID(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "User nonexistent not found")
}

func TestGetAllUsers_Empty(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	repo.On("GetAllUsers", mock.Anything).Return([]*domain.User{}, nil)

	users, err := s.GetAllUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}

func TestGetAllUsers_Failure(t *testing.T) {
	repo := new(MockUserRepository)
	s := service.NewUserService(repo)

	repo.On("GetAllUsers", mock.Anything).Return([]*domain.User{}, errors.New("db error"))

	_, err := s.GetAllUsers(context.Background())
	assert.Error(t, err)
}