package repository

import (
	"telegram-go-bot/internal/application/model"
	"telegram-go-bot/internal/service"
)

type RedisAdapter struct {
	redis *service.RedisIntegration
}

func NewRedisAdapter(redis *service.RedisIntegration) *RedisAdapter {
	return &RedisAdapter{
		redis: redis,
	}
}

func (adp *RedisAdapter) SaveAccessToken(accessToken *model.AccessToken, userId string) error {
	return adp.redis.SaveAccessToken(accessToken, userId)
}

func (adp *RedisAdapter) RetrieveAccessToken(userId string) (*model.AccessToken, error) {
	return adp.redis.RetrieveAccessToken(userId)
}

func (adp *RedisAdapter) IncreaseUserLoginAttempt(userId string) error {
	return adp.redis.IncreaseUserLoginAttempt(userId)
}

func (adp *RedisAdapter) RetrieveUserLoginAttempt(userId string) (int, error) {
	return adp.redis.RetrieveUserLoginAttempt(userId)
}

func (adp *RedisAdapter) ResetUserLoginAttemp(userId string) error {
	return adp.redis.ResetUserLoginAttemp(userId)
}
