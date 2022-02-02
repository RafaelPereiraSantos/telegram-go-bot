package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"telegram-go-bot/internal/adapter/in"
	apiOut "telegram-go-bot/internal/adapter/out/api"
	"telegram-go-bot/internal/adapter/out/repository"
	"telegram-go-bot/internal/api"
	"telegram-go-bot/internal/application/service"
	"telegram-go-bot/internal/bot"
	extService "telegram-go-bot/internal/service"
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
	reddit := extService.NewRedditIntegration(
		os.Getenv("REDDIT_APP_ID"),
		os.Getenv("REDDIT_APP_TOKEN"),
	)
	redditAdp := apiOut.NewReddtAdp(reddit)

	redis, err := extService.NewRedisIntegration(
		os.Getenv("REDIS_ADDRESS"),
		os.Getenv("REDIS_PASS"),
		0,
	)

	if err != nil {
		msg := fmt.Sprintf("Unable to connect with redis %s\n", err.Error())
		fmt.Println(msg)
		return
	}

	redisAdp := repository.NewRedisAdapter(redis)

	botService := service.NewBot(redditAdp, redisAdp)
	botServiceAdp := in.NewBotImpl(botService)

	b, err := bot.NewTelegramBot(os.Getenv("TELEGRAM_SECRET_TOKEN"), botServiceAdp)
	if err != nil {
		msg := fmt.Sprintf("Unable to start bot %s\n", err.Error())
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
