package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	startBot()

	startHealth()
}

func startBot() {
	telegramToken := os.Getenv("TELEGRAM_SECRET_TOKEN")

	bot, err := botApi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botApi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := botApi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func startHealth() {
	port := os.Getenv("PORT")

	fmt.Printf("started server at :%s", port)

	http.HandleFunc("/health", handler)
	http.ListenAndServe(":"+port, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("teste")
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
