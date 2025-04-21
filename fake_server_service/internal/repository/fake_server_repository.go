package repository

import (
	"context"
	
	"github.com/redis/go-redis/v9"
)

type FakeServerRepository interface {
	EnableServer(serverID int) error
	CountEnabledServers() (int, error)
	GetServerStatus(serverID int) (bool, error)
}

type fakeServerRepository struct {
	redisClient *redis.Client
}

func NewFakeServerRepository(redisClient *redis.Client) FakeServerRepository {
	return &fakeServerRepository{
		redisClient: redisClient,
	}
}

func (r *fakeServerRepository) EnableServer(serverID int) error {
	cmd := r.redisClient.SetBit(context.Background(), "enabled_servers", int64(serverID), 1)
	return cmd.Err()
}

func (r *fakeServerRepository) CountEnabledServers() (int, error) {
	cmd := r.redisClient.BitCount(context.Background(), "enabled_servers", &redis.BitCount{
		Start: 0,
		End:   -1,
	})

	if err := cmd.Err(); err != nil {
		return 0, err
	}

	count, err := cmd.Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *fakeServerRepository) GetServerStatus(serverID int) (bool, error) {
	cmd := r.redisClient.GetBit(context.Background(), "enabled_servers", int64(serverID))
	if err := cmd.Err(); err != nil {
		return false, err
	}

	status, err := cmd.Result()
	if err != nil {
		return false, err
	}

	return status == 1, nil
}