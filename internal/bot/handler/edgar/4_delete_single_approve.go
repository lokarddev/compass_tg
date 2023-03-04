package edgar

import (
	"app/internal/bot/handler"
	"app/internal/bot/helper"
	"app/internal/repository"
	"fmt"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DelSingleApproveHandler struct {
	bot  *api.BotAPI
	repo repository.BaseRepoInterface
	handler.BaseHandler
	nav handler.Nav
}

func NewDelSingleApproveHandler(bot *api.BotAPI, dbRepo repository.BaseRepoInterface, branch int) *DelSingleApproveHandler {
	return &DelSingleApproveHandler{
		bot:  bot,
		repo: dbRepo,
		nav: handler.Nav{
			ValidSources: []int{helper.EdgarDeleteSingleState},
			ValidBranch:  branch,
		}}
}

func (h *DelSingleApproveHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user, h.nav) {
		return
	}

	toSub := make(map[string]string)

	for _, ticker := range user.Subscriptions.Edgar.PendingSubs {
		if company, ok := handler.MockEdgarSubs[ticker]; ok {
			toSub[ticker] = company
		}
	}

	msgText := helper.MsgDelSingle
	for ticker, company := range toSub {
		msgText += fmt.Sprintf("%s - %s\n", ticker, company)
	}

	msg := api.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyMarkup = helper.SubApproveKeyboard

	if _, err = h.bot.Send(msg); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, helper.EdgarDelSingleApproveState, helper.EdgarBranch); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
