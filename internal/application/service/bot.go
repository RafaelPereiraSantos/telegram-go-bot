package service

import (
	"fmt"

	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/model"
)

type Bot struct{}

func NewBot() *Bot {
	return &Bot{}
}

func (bot *Bot) ReceiveMessage(message model.Message) string {
	msg := fmt.Sprintf("Message [%s] Received From [%s]", message.Text, message.User)
	fmt.Println(msg)

	return "stop bothering me"
}
