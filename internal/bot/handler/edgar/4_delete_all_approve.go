package edgar

import (
	"app/internal/bot/handler"
	"app/internal/bot/helper"
	"app/internal/repository"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var edgarDelAllApproveButtons = []string{helper.DeleteAllButton}

type DelAllApproveHandler struct {
	bot  *api.BotAPI
	repo repository.BaseRepoInterface
	handler.BaseHandler
	nav handler.Nav
}

func NewDelAllApproveHandler(bot *api.BotAPI, dbRepo repository.BaseRepoInterface, branch int) *DelAllApproveHandler {
	return &DelAllApproveHandler{
		bot:  bot,
		repo: dbRepo,
		nav: handler.Nav{
			ValidSources: []int{helper.EdgarDeleteAllState},
			ValidBranch:  branch,
		}}
}

func (h *DelAllApproveHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user, h.nav) || !h.ValidText(update.Message.Text, edgarDelAllApproveButtons) {
		return
	}

	msgText := helper.MsgDelAll

	msg := api.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyMarkup = helper.SubApproveKeyboard

	if _, err = h.bot.Send(msg); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, helper.EdgarDelAllApproveState, helper.EdgarBranch); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
