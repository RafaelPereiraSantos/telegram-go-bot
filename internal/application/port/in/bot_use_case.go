package in

import "telegram-go-bot/internal/application/model"

type BotUseCase interface {
	ReceiveMessage(message model.ReceivedMessage) []model.ReplyMessage
}
