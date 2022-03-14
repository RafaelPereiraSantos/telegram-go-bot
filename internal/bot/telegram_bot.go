package bot

import (
	"log"

	"telegram-go-bot/internal/application/model"
	portIn "telegram-go-bot/internal/application/port/in"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot           *telegramBot.BotAPI
	messageSevice portIn.BotUseCase
}

func NewTelegramBot(telegramToken string, messageSevice portIn.BotUseCase) (*TelegramBot, error) {
	bot, err := telegramBot.NewBotAPI(telegramToken)

	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:           bot,
		messageSevice: messageSevice,
	}, nil
}

// ListenEvents, it starts the event listener of the telegram bot api.
func (tBot *TelegramBot) ListenEvents(debug bool) error {

	bot := tBot.bot
	bot.Debug = debug

	srv := tBot.messageSevice

	updateConfig := telegramBot.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		receivedMessage := update.Message

		log.Printf("[%s] %s\n", receivedMessage.From.UserName, receivedMessage.Text)

		received := model.ReceivedMessage{
			ChatId: update.Message.Chat.ID,
			User:   receivedMessage.From.UserName,
			Text:   receivedMessage.Text,
		}

		replyMessages := srv.ReceiveMessage(received)

		chatID := receivedMessage.Chat.ID
		for _, reply := range replyMessages {
			if reply.Image == nil {
				bot.Send(telegramBot.NewMessage(chatID, reply.Text))
			} else {
				msg := telegramBot.NewDocument(chatID, reply.Image)
				msg.Caption = reply.Text
				bot.Send(msg)
			}
		}
	}

	return nil
}
