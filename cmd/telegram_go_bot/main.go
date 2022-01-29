package main

import (
	"fmt"
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
		fmt.Println("Error loading .env file")
	}

	go startBot()

	startHealth()
}

func startBot() {
	botService := service.NewBot()
	botServiceAdp := in.NewBotImpl(botService)

	b, err := bot.NewTelegramBot(os.Getenv("TELEGRAM_SECRET_TOKEN"), botServiceAdp)
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
