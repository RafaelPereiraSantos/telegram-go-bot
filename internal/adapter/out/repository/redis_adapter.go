package repository

import "telegram-go-bot/internal/service"

type RedisAdapter struct {
	redis service.RedisIntegration
}

func NewRedisAdapter(redis service.RedisIntegration) *RedisAdapter {
	return &RedisAdapter{
		redis: redis,
	}
}

func (adp *RedisAdapter) SaveUserAsAdmin(adminAccount string) error {
	return adp.redis.SaveUserAsAdmin(adminAccount)
}

func (adp *RedisAdapter) IsAdminUser(account string) (bool, error) {
	return adp.redis.IsAdminUser(account)
}

func (adp *RedisAdapter) IncreaseUserLoginAttempt(account string) error {
	return adp.redis.IncreaseUserLoginAttempt(account)
}

func (adp *RedisAdapter) ResetUserLoginAttemp(account string) error {
	return adp.redis.ResetUserLoginAttemp(account)
}
