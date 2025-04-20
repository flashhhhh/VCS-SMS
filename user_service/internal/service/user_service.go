package service

import (
	"context"
	"errors"
	"time"
	"user_service/internal/domain"
	"user_service/internal/repository"

	"github.com/flashhhhh/pkg/hash"
	"github.com/flashhhhh/pkg/jwt"
	"github.com/flashhhhh/pkg/logging"
)

type UserService interface {
	CreateUser(ctx context.Context, username, password, name, email, role string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	logging.LogMessage("user_service", "Initializing UserService", "INFO")

	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, password, name, email, role string) (string, error) {
	hashedPassword := hash.HashString(password)

	user := &domain.User{
		Username: username,
		Password: hashedPassword,
		Name:     name,
		Email:    email,
		Role:     role,
	}

	userID, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepository.Login(ctx, username)
	if err != nil {
		return "", err
	}

	if !hash.CompareHashAndString(user.Password, password) {
		return "", errors.New("Invalid password")
	}

	token, err := jwt.GenerateToken(
		map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		}, time.Hour)
	
	if err != nil {
		return "", errors.New("Failed to generate token: " + err.Error())
	}
	return token, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("User " + id + " not found: " + err.Error())
	}
	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, errors.New("Failed to retrieve all users: " + err.Error())
	}
	return users, nil
}