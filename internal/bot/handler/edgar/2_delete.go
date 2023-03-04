package edgar

import (
	"app/internal/bot/handler"
	"app/internal/bot/helper"
	"app/internal/repository"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeleteHandler struct {
	bot  *api.BotAPI
	repo repository.BaseRepoInterface
	handler.BaseHandler
	nav handler.Nav
}

var (
	edgarDeleteButtons = []string{helper.DeleteButton}
)

func NewDeleteHandler(bot *api.BotAPI, dbRepo repository.BaseRepoInterface, branch int) *DeleteHandler {
	return &DeleteHandler{
		bot:  bot,
		repo: dbRepo,
		nav: handler.Nav{
			ValidSources: []int{helper.EdgarCommandState},
			ValidBranch:  branch,
		}}
}

func (h *DeleteHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user, h.nav) || !h.ValidText(update.Message.Text, edgarDeleteButtons) {
		return
	}

	msg := api.NewMessage(update.Message.Chat.ID, helper.MsgTypeTickersToUnsubscribe)
	msg.ReplyMarkup = helper.DeleteAllSubKeyboard

	if _, err = h.bot.Send(msg); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, helper.EdgarDeleteState, helper.EdgarBranch); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
