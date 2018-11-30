package notifications

import (
	"reviewer/api/config"
	"reviewer/api/store"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

var (
	bot *tgbotapi.BotAPI
)

func init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(config.TgAPIKey)
	if err != nil {
		logrus.Errorf("Cannot create bot api: %+v", err)
		bot = nil
	}
}

// TelegramSend sends notification to users telegram
func TelegramSend(user store.User, message string) {
	if bot == nil {
		return
	}
	if user.TelegramID <= 0 {
		return
	}
	go func() {
		msg := tgbotapi.NewMessage(int64(user.TelegramID), message)
		msg.ParseMode = tgbotapi.ModeMarkdown
		_, err := bot.Send(msg)
		if err != nil {
			logrus.Errorf("Cannot send notification: %+v", err)
			return
		}
		return
	}()
}
