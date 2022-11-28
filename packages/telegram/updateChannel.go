package telegram

import (
	"pinger/packages/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewUpdateChannel(
	bot *tgbotapi.BotAPI,
	config config.Config,
) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = int(config.Bot.UpdateTimeout)

	var channel tgbotapi.UpdatesChannel

	if config.Bot.ListenerEnabled {
		channel = bot.GetUpdatesChan(u)
	} else {
		channel = *new(tgbotapi.UpdatesChannel)
	}

	return channel
}
