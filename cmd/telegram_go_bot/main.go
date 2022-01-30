package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/adapter/in"
	apiOut "github.com/RafaelPereiraSantos/telegram-go-bot/internal/adapter/out/api"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/api"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/service"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/bot"
	extService "github.com/RafaelPereiraSantos/telegram-go-bot/internal/service"
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
	user := extService.RedditUser{
		UserName: os.Getenv("REDDIT_USER_NAME"),
		Password: os.Getenv("REDDIT_USER_PASSWORD"),
	}

	reddit := extService.NewRedditIntegration(
		os.Getenv("REDDIT_APP_ID"),
		os.Getenv("REDDIT_APP_TOKEN"),
		user,
	)
	redditAdp := apiOut.NewReddtAdp(reddit)

	botService := service.NewBot(redditAdp)
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
