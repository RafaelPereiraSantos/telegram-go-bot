package out

import "telegram-go-bot/internal/application/model"

type UserRepository interface {
	SaveAccessToken(accessToken *model.AccessToken, userId string) error
	RetrieveAccessToken(userId string) (*model.AccessToken, error)
	IncreaseUserLoginAttempt(userId string) error
	RetrieveUserLoginAttempt(userId string) (int, error)
	ResetUserLoginAttemp(userId string) error
}
