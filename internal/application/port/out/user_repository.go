package out

import "telegram-go-bot/internal/application/model"

type UserRepository interface {
	SaveAccessToken(accessToken *model.AccessToken, account string) error
	RetrieveAccessToken(account string) (*model.AccessToken, error)
	IncreaseUserLoginAttempt(account string) error
	RetrieveUserLoginAttempt(account string) (int, error)
	ResetUserLoginAttemp(account string) error
}
