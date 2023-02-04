package common

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Edgar keyboards
var (
	DeleteSubKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Delete"),
		))
)
