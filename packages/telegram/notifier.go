package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Notifier interface {
	Notify(message string, target int64) error
}

func NewNotifier(
	logger *zap.SugaredLogger,
	bot *tgbotapi.BotAPI,
) Notifier {
	return &notifier{
		logger: logger,
		bot:    bot,
	}
}

type notifier struct {
	logger *zap.SugaredLogger
	bot    *tgbotapi.BotAPI
}

func (n *notifier) Notify(message string, target int64) error {
	msg := tgbotapi.NewMessage(target, message)

	_, err := n.bot.Send(msg)

	return err
}
