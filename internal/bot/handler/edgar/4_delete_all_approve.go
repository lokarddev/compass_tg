package edgar

import (
	"app/internal/bot/common/keyboards"
	"app/internal/bot/common/message"
	"app/internal/bot/handler"
	"app/internal/bot/repository/mongo_db"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	edgarDelAllApproveState = "EdgarDelAllApprove"
)

var (
	edgarDelAllApproveSrc     = []string{edgarDeleteAllState}
	edgarDelAllApproveButtons = []string{keyboards.DeleteAllButton}
)

type DelAllApproveHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
	handler.BaseHandler
}

func NewDelAllApproveHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *DelAllApproveHandler {
	return &DelAllApproveHandler{bot: bot, repo: dbRepo}
}

func (h *DelAllApproveHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user.State.NavCurrent, edgarDelAllApproveSrc) || !h.ValidText(update.Message.Text, edgarDelAllApproveButtons) {
		return
	}

	msgText := message.MsgDelAll

	msg := api.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyMarkup = keyboards.SubApproveKeyboard

	if _, err = h.bot.Send(msg); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, edgarDelAllApproveState); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
