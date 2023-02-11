package edgar

import (
	"app/internal/bot/common"
	"app/internal/bot/common/keyboards"
	"app/internal/bot/common/message"
	"app/internal/bot/handler"
	"app/internal/bot/repository/mongo_db"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type DeleteHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
	handler.BaseHandler
}

const (
	edgarDeleteState = "edgarDeleteState"
)

var (
	edgarDeleteSources = []string{common.EdgarCommand}
	edgarDeleteButtons = []string{keyboards.DeleteButton}
)

func NewDeleteHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *DeleteHandler {
	return &DeleteHandler{bot: bot, repo: dbRepo}
}

func (h *DeleteHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user.State.NavCurrent, edgarDeleteSources) || !h.ValidText(update.Message.Text, edgarDeleteButtons) {
		return
	}

	msg := api.NewMessage(update.Message.Chat.ID, message.MsgTypeTickersToUnsubscribe)
	msg.ReplyMarkup = keyboards.DeleteAllSubKeyboard

	if _, err = h.bot.Send(msg); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, edgarDeleteState); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
