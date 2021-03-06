package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"telegram-go-bot/internal/application/port/in"
)

type (
	HealthApi struct {
		service in.HealthCheckUseCase
	}
)

func NewHealthApi(service in.HealthCheckUseCase) *HealthApi {
	return &HealthApi{
		service: service,
	}
}

func (api *HealthApi) Start(port string) {
	fmt.Printf("started server at :%s\n", port)

	http.HandleFunc("/health", api.healthStatus)
	http.ListenAndServe(":"+port, nil)
}

func (api *HealthApi) healthStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("checking health...\n")

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(api.service.Heatlh())
}
