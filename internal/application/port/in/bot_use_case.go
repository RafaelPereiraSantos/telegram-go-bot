package in

import "github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/model"

type BotUseCase interface {
	ReceiveMessage(message model.Message) string
}
