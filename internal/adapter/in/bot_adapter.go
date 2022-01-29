package in

import (
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/model"
	"github.com/RafaelPereiraSantos/telegram-go-bot/internal/application/service"
)

type (
	BotAdp interface {
		ReceiveMessage(userName, message string) string
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

func (impl *BotImpl) ReceiveMessage(userName, message string) string {
	msg := model.Message{
		User: userName,
		Text: message,
	}

	return impl.srv.ReceiveMessage(msg)
}
