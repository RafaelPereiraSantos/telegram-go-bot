package service

type CheckServicesHealth struct{}

func NewCheckServicesHealth() *CheckServicesHealth {
	return &CheckServicesHealth{}
}

func (serice *CheckServicesHealth) Heatlh() []map[string]string {
	var status []map[string]string

	applicationStatus := map[string]string{
		"status": "OK",
	}

	status = append(status, applicationStatus)

	return status
}
