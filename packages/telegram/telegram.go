package telegram

import (
	"context"
	"pinger/packages/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewBot(
	lc fx.Lifecycle,
	config config.Config,
	logger *zap.SugaredLogger,
) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(config.Bot.Token)

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			logger.Infof("Authorized on account %s", bot.Self.UserName)

			bot.Debug = config.Bot.Debug
			_, err := bot.Request(tgbotapi.DeleteWebhookConfig{
				DropPendingUpdates: true,
			})

			return err
		},
	})

	return bot, err
}
