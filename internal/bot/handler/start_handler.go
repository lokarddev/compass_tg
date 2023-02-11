package handler

import (
	"app/internal/bot/common"
	"app/internal/bot/common/message"
	"app/internal/bot/repository/mongo_db"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
}

func NewStartHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *StartHandler {
	return &StartHandler{bot: bot, repo: dbRepo}
}

func (h *StartHandler) Call(update *api.Update) {
	user, err := h.repo.UpsertUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot upsert user, due to error: %s", err.Error())

		return
	}

	msg := api.NewMessage(update.Message.Chat.ID, "")
	msg.Text = message.MsgStart
	msg.ReplyMarkup = api.NewRemoveKeyboard(false)

	if _, err = h.bot.Send(msg); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, common.StartCommand); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
