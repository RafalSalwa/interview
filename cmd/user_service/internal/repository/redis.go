package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/models"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	log         *logger.Logger
	redisClient redis.UniversalClient
}

func (r Redis) PutUser(ctx context.Context, user *models.UserDBModel) error {
	key := fmt.Sprintf("user_%d", user.Id)
	redisUser := &models.UserRedisModel{}
	redisUser.FromDbModel(user)
	bytes, err := json.Marshal(redisUser)

	if err != nil {
		r.log.Error().Err(err).Msg("redis:user:put:marshal")
		return err
	}

	if errR := r.redisClient.HSetNX(ctx, "auth:repo", key, bytes).Err(); errR != nil {
		r.log.Error().Err(errR).Msg("redis:user:put:HSetNX")
		return errR
	}
	return nil
}

func NewRedisRepository(redisClient redis.UniversalClient, log *logger.Logger) *Redis {
	return &Redis{log: log, redisClient: redisClient}
}
