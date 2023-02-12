package edgar

import (
	"app/internal/bot/common/keyboards"
	"app/internal/bot/common/message"
	"app/internal/bot/handler"
	"app/internal/bot/repository/mongo_db"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	edgarDelSingleApproveState = "EdgarDelSingleApprove"
)

var (
	edgarDelSingleApproveSrc = []string{edgarDeleteSingleState}
)

type DelSingleApproveHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
	handler.BaseHandler
}

func NewDelSingleApproveHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *DelSingleApproveHandler {
	return &DelSingleApproveHandler{bot: bot, repo: dbRepo}
}

func (h *DelSingleApproveHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user.State.NavCurrent, edgarDelSingleApproveSrc) {
		return
	}

	toSub := make(map[string]string)

	for _, ticker := range user.Subscriptions.Edgar.PendingSubs {
		if company, ok := handler.MockEdgarSubs[ticker]; ok {
			toSub[ticker] = company
		}
	}

	msgText := message.MsgDelSingle
	for ticker, company := range toSub {
		msgText += fmt.Sprintf("%s - %s\n", ticker, company)
	}

	msg := api.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyMarkup = keyboards.SubApproveKeyboard

	if _, err = h.bot.Send(msg); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, edgarDelSingleApproveState); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
