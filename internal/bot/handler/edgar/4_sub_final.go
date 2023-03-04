package edgar

import (
	"app/internal/bot/handler"
	"app/internal/bot/helper"
	"app/internal/model"
	"app/internal/repository"
	"fmt"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	edgarSubFinalButtons = []string{helper.Yes, helper.No}
)

type SubFinalHandler struct {
	bot  *api.BotAPI
	repo repository.BaseRepoInterface
	handler.BaseHandler
	nav handler.Nav
}

func NewSubFinalHandler(bot *api.BotAPI, dbRepo repository.BaseRepoInterface, branch int) *SubFinalHandler {
	return &SubFinalHandler{
		bot:  bot,
		repo: dbRepo,
		nav: handler.Nav{
			ValidSources: []int{helper.EdgarSubscribeApproveState},
			ValidBranch:  branch,
		}}
}

func (h *SubFinalHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user, h.nav) || !h.ValidText(update.Message.Text, edgarSubFinalButtons) {
		return
	}

	msg := api.NewMessage(update.Message.Chat.ID, helper.MsgSubSuccess)
	msg.ReplyMarkup = api.NewRemoveKeyboard(false)

	switch update.Message.Text {
	case helper.No:
		msg.Text = helper.MsgNoSub

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

	case helper.Yes:
		toSub := user.Subscriptions.Edgar.PendingSubs

		if err = h.repo.UpsertSubscriptions(&user, model.EdgarSubscription); err != nil {
			log.Println(err)

			return
		}

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

		msg = api.NewMessage(update.Message.Chat.ID, "")
		msg.Text = mockLastTickers(toSub)

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}
	}

	if err = h.repo.UpsertState(&user, helper.StartCommandState, helper.StartBranch); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}

func mockLastTickers(subs []string) string {
	msg := "Here are the last filings\n"
	for _, ticker := range subs {
		msg += fmt.Sprintf("%s\nfucking first filing\nfucking second filing\n", ticker)
	}

	return msg
}
