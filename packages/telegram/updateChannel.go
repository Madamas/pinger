package telegram

import (
	"pinger/packages/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func NewUpdateChannel(
	bot *tgbotapi.BotAPI,
	config config.Config,
	logger *zap.SugaredLogger,
) tgbotapi.UpdatesChannel {
	var channel tgbotapi.UpdatesChannel

	if config.Bot.ListenerEnabled {
		logger.Info("Starting bot with telegram listener enabled")

		u := tgbotapi.NewUpdate(0)
		u.Timeout = int(config.Bot.UpdateTimeout)
		channel = bot.GetUpdatesChan(u)
	} else {
		logger.Info("Starting bot with telegram listener disabled")
		channel = *new(tgbotapi.UpdatesChannel)
	}

	return channel
}
