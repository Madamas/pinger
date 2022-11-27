package telegram

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"pinger/packages/config"
	"pinger/packages/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Listener interface {
	Handle()
}

func NewListener(
	logger *zap.SugaredLogger,
	bot *tgbotapi.BotAPI,
	channel tgbotapi.UpdatesChannel,
	storage storage.Storage,
	client *http.Client,
	config config.Config,
	lc fx.Lifecycle,
) Listener {
	listener := listener{
		logger:  logger,
		bot:     bot,
		channel: channel,
		storage: storage,
		client:  client,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if config.Bot.ListenerEnabled {
					go listener.Handle()
				}

				return nil
			},
		},
	)

	return &listener
}

type listener struct {
	logger  *zap.SugaredLogger
	bot     *tgbotapi.BotAPI
	channel tgbotapi.UpdatesChannel
	storage storage.Storage
	client  *http.Client
}

func (l *listener) Handle() {
	for update := range l.channel {
		if update.Message != nil {
			receivedText := update.Message.Text
			response := "empty"

			_, err := l.storage.AddUser(update.Message.From.ID)

			if err != nil {
				l.logger.Errorf("Error while creating user, err - %s", err.Error())
			}

			status, err := l.checkStatus(update.Message.From.ID)

			if err != nil {
				l.logger.Errorf("Error while fetching user status, err - %s", err.Error())
			}

			command := ParseCommand(strings.TrimSpace(receivedText))

			if command.IsOk() {
				if command == COMMAND_START || command == COMMAND_BACK {
					// new status, err
					l.resetStatus(update.Message.From.ID)
					status = storage.STATUS_INITIAL
				}

				if command == COMMAND_NEW_TARGET {
					// err
					l.storage.SetStatus(update.Message.From.ID, storage.STATUS_NEW_TARGET)
					status = storage.STATUS_NEW_TARGET
				}

				if command == COMMAND_LIST_TARGETS {
					// err
					if targets, err := l.storage.FetchTargetsByOwner(update.Message.From.ID); err != nil {
						l.logger.Errorf("Couldn't fetch targets for user %d, error - %s", update.Message.From.ID, err.Error())
						response = "Something went wrong :("
					} else {
						response = responseWithTargets(targets)
					}
				}
			} else {
				if status == storage.STATUS_NEW_TARGET {
					// response
					response = l.addTarget(receivedText, update.Message.From.ID)
					l.resetStatus(update.Message.From.ID)
					status = storage.STATUS_INITIAL
				}
			}

			l.logger.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			// msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = NewKeyboard(status)

			_, err = l.bot.Send(msg)
			if err != nil {
				l.logger.Errorf("Something went wrong %s", err.Error())
			}
		}
	}
}

func responseWithTargets(targets []storage.Target) string {
	response := make([]string, 0, len(targets))

	for idx, target := range targets {
		response = append(response, fmt.Sprintf("%d. %s", idx+1, target.Url))
	}

	return strings.Join(response, "\n")
}

func (l *listener) addTarget(text string, userId int64) string {
	_, err := url.Parse(text)

	if err != nil {
		return fmt.Sprintf("Incorrect URL: %s", text)
	}

	r, err := l.client.Get(text)

	if err != nil || r.StatusCode > 300 {
		var eerror string
		if err == nil {
			eerror = fmt.Sprintf("Endpoint returned status %d", r.StatusCode)
		} else {
			eerror = err.Error()
		}

		l.logger.Errorf("Invalid endpoint: %s", eerror)

		return "Something went wrong :("
	}

	if err := l.storage.NewTarget(userId, text); err != nil {
		l.logger.Errorf("Unable to add new target %s", err.Error())
	}

	return "Ok"
}

func (l *listener) resetStatus(userId int64) error {
	if err := l.storage.SetStatus(userId, storage.STATUS_INITIAL); err != nil {
		return err
	}

	return nil
}

func (l *listener) checkStatus(userId int64) (storage.Status, error) {
	return l.storage.GetStatus(userId)
}
