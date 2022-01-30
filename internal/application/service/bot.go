package service

import (
	"fmt"

	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/model"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/port/out/api"
)

type Bot struct {
	scraper api.SocialMediaScraper
}

func NewBot(scraper api.SocialMediaScraper) *Bot {
	return &Bot{
		scraper: scraper,
	}
}

func (bot *Bot) ReceiveMessage(message model.Message) string {
	msg := fmt.Sprintf("Message [%s] Received From [%s]", message.Text, message.User)
	fmt.Println(msg)

	return "stop bothering me"
}
