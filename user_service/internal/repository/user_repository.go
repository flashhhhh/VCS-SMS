package repository

import (
	"context"
	"user_service/internal/domain"

	"github.com/flashhhhh/pkg/logging"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	Login(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
	redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepository {
	logging.LogMessage("user_service", "Initializing UserRepository", "INFO")

	return &userRepository{
		db: db,
		redis: redis,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	// Generate a unique ID for the user
	user.ID = uuid.New().String()

	err := r.db.Create(user).Error
	if err != nil {
		return "", err
	}
	
	return user.ID, nil
}

func (r *userRepository) Login(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}