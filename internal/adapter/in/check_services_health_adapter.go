package in

import "telegram-go-bot/internal/application/service"

type (
	HealthCheckAdp interface{}

	HealthCheckImp struct {
		srv *service.CheckServicesHealth
	}
)

func NewHealthCheckImp(srv *service.CheckServicesHealth) *HealthCheckImp {
	return &HealthCheckImp{
		srv: srv,
	}
}

func (impl *HealthCheckImp) Heatlh() []map[string]string {
	return impl.srv.Heatlh()
}
