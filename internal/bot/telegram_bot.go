package bot

import (
	"log"

	telegramBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot *telegramBot.BotAPI
}

func NewTelegramBot(telegramToken string) (*TelegramBot, error) {
	bot, err := telegramBot.NewBotAPI(telegramToken)

	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot: bot,
	}, nil
}

func (tBot *TelegramBot) ListenEvents(debug bool) error {

	bot := tBot.bot
	bot.Debug = debug

	updateConfig := telegramBot.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := telegramBot.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

	return nil
}
