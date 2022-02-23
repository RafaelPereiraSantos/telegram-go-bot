package in

import (
	"telegram-go-bot/internal/application/model"
	"telegram-go-bot/internal/application/service"
)

type (
	BotAdp interface {
		ReceiveMessage(chatId int64, userName, message string) string
	}

	BotImpl struct {
		srv *service.Bot
	}
)

func NewBotImpl(srv *service.Bot) *BotImpl {
	return &BotImpl{
		srv: srv,
	}
}

func (impl *BotImpl) ReceiveMessage(msg model.ReceivedMessage) []model.ReplyMessage {
	return impl.srv.ReceiveMessage(msg)
}
