package repository

import (
	"context"

	"github.com/RafalSalwa/auth-api/pkg/models"
	"github.com/go-redis/redis/v8"
)

type RedisAdapter struct {
	DB *redis.Client
}

func (r RedisAdapter) Exists(ctx context.Context, udb *models.UserDBModel) bool {
	_, err := r.DB.Exists(ctx, udb.Email).Result()
	return err == nil
}

func newRedisUserRepository(client *redis.Client) UserRepository {
	return &RedisAdapter{DB: client}
}

func (r RedisAdapter) Update(ctx context.Context, user *models.UserDBModel) error {
	_, err := r.DB.Set(ctx, user.Email, user, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisAdapter) Save(ctx context.Context, user *models.UserDBModel) error {
	_, err := r.DB.Set(ctx, user.Email, user, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisAdapter) FindAll(ctx context.Context, user *models.UserDBModel) ([]models.UserDBModel, error) {
	panic("implement me")
}

func (r RedisAdapter) FindOne(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	_, err := r.DB.Get(ctx, user.Email).Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r RedisAdapter) GetOrCreate(ctx context.Context, id int64) (*models.UserDBModel, error) {
	// TODO implement me
	panic("implement me")
}

func (r RedisAdapter) Confirm(ctx context.Context, udb *models.UserDBModel) error {
	_, err := r.DB.Set(ctx, udb.Email, udb, 0).Result()
	if err != nil {
		return err
	}
	return nil
}
