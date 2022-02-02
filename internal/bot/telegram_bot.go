package bot

import (
	"log"

	"telegram-go-bot/internal/adapter/in"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot           *telegramBot.BotAPI
	messageSevice in.BotAdp
}

func NewTelegramBot(telegramToken string, messageSevice in.BotAdp) (*TelegramBot, error) {
	bot, err := telegramBot.NewBotAPI(telegramToken)

	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:           bot,
		messageSevice: messageSevice,
	}, nil
}

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

		log.Printf("[%s] %s", receivedMessage.From.UserName, receivedMessage.Text)

		replyText := srv.ReceiveMessage(
			update.Message.Chat.ID,
			receivedMessage.From.UserName,
			receivedMessage.Text,
		)

		msg := telegramBot.NewMessage(receivedMessage.Chat.ID, replyText)
		// msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

	return nil
}
