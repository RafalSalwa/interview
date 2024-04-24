package repository

import (
	"context"

	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	log         *logger.Logger
	redisClient redis.UniversalClient
}

type RedisAdapter struct {
	DB *redis.UniversalClient
}

func (r RedisAdapter) Exists(ctx context.Context, udb *models.UserDBModel) bool {
	panic("implement me, Redis Exists")
}

func newRedisUserRepository(client *redis.UniversalClient) UserRepository {
	return &RedisAdapter{DB: client}
}

func (r RedisAdapter) Update(ctx context.Context, user models.UserDBModel) error {
	// TODO implement me
	panic("implement me")
}

func (r RedisAdapter) Save(ctx context.Context, user *models.UserDBModel) error {
	// TODO implement me
	panic("implement me")
}

func (r RedisAdapter) FindAll(ctx context.Context, user *models.UserDBModel) ([]models.UserDBModel, error) {
	// TODO implement me
	panic("implement me")
}

func (r RedisAdapter) FindOne(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	// TODO implement me
	panic("implement me")
}

func (r RedisAdapter) GetOrCreate(ctx context.Context, id int64) (*models.UserDBModel, error) {
	// TODO implement me
	panic("implement me")
}

func (r RedisAdapter) Confirm(ctx context.Context, udb *models.UserDBModel) error {
	// TODO implement me
	panic("implement me")
}
