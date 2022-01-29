package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/adapter/in"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/api"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/service"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// startBot()

	startHealth()
}

// func startBot() {
// 	telegramToken := os.Getenv("TELEGRAM_SECRET_TOKEN")

// 	bot, err := botApi.NewBotAPI(telegramToken)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	bot.Debug = true

// 	log.Printf("Authorized on account %s", bot.Self.UserName)

// 	u := botApi.NewUpdate(0)
// 	u.Timeout = 60

// 	updates, err := bot.GetUpdatesChan(u)

// 	for update := range updates {
// 		if update.Message == nil {
// 			continue
// 		}

// 		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

// 		msg := botApi.NewMessage(update.Message.Chat.ID, update.Message.Text)
// 		msg.ReplyToMessageID = update.Message.MessageID

// 		bot.Send(msg)
// 	}
// }

func startHealth() {
	healthService := service.NewCheckServicesHealth()
	healthServiceAdp := in.NewHealthCheckImp(healthService)

	healthApi := api.NewHealthApi(healthServiceAdp)
	healthApi.Start(os.Getenv("PORT"))
}
