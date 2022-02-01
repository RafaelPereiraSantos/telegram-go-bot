package out

type UserRepository interface {
	SaveUserAsAdmin(adminAccount string) error
	IsAdminUser(account string) (bool, error)
	IncreaseUserLoginAttempt(account string) error
	ResetUserLoginAttemp(account string) error
}
