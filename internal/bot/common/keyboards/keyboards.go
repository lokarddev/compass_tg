package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Buttons
const (
	DeleteButton    = "Delete"
	DeleteAllButton = "Delete All"
	Yes             = "Yes"
	No              = "No"
)

// Edgar keyboards
var (
	DeleteSubKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DeleteButton),
		))

	DeleteAllSubKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DeleteAllButton),
		))

	SubApproveKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(Yes),
			tgbotapi.NewKeyboardButton(No),
		))
)
