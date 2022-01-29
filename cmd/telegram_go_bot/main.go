package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/adapter/in"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/api"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/service"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/bot"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	go startBot()

	startHealth()
}

func startBot() {
	b, err := bot.NewTelegramBot(os.Getenv("TELEGRAM_SECRET_TOKEN"))

	if err != nil {
		msg := fmt.Sprintf("Unable to start bot %s", err.Error())
		fmt.Println(msg)
		return
	}

	b.ListenEvents(true)
}

func startHealth() {
	healthService := service.NewCheckServicesHealth()
	healthServiceAdp := in.NewHealthCheckImp(healthService)

	healthApi := api.NewHealthApi(healthServiceAdp)
	healthApi.Start(os.Getenv("PORT"))
}
