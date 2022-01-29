package service

import "github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/model"

type Bot struct{}

func NewBot() *Bot {
	return &Bot{}
}

func (bot *Bot) ReceiveMessage(message model.Message) string {
	return "stop bothering me"
}
