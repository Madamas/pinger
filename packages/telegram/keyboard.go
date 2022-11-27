package telegram

import (
	"pinger/packages/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	TYPICAL_KEYBOARD = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(COMMAND_NEW_TARGET.String()),
			tgbotapi.NewKeyboardButton(COMMAND_LIST_TARGETS.String()),
		),
	)

	EDIT_KEYBOARD = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(COMMAND_BACK.String()),
		),
	)
)

func NewKeyboard(status storage.Status) tgbotapi.ReplyKeyboardMarkup {
	if status == storage.STATUS_INITIAL {
		return TYPICAL_KEYBOARD
	}

	return EDIT_KEYBOARD
}
