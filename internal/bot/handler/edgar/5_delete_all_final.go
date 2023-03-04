package edgar

import (
	"app/internal/bot/handler"
	"app/internal/bot/helper"
	"app/internal/model"
	"app/internal/repository"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DelAllFinalHandler struct {
	bot  *api.BotAPI
	repo repository.BaseRepoInterface
	handler.BaseHandler
	nav handler.Nav
}

func NewDelAllFinalHandler(bot *api.BotAPI, dbRepo repository.BaseRepoInterface, branch int) *DelAllFinalHandler {
	return &DelAllFinalHandler{
		bot:  bot,
		repo: dbRepo,
		nav: handler.Nav{
			ValidSources: []int{helper.EdgarDeleteAllState},
			ValidBranch:  branch,
		}}
}

func (h *DelAllFinalHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user, h.nav) || !h.ValidText(update.Message.Text, edgarSubFinalButtons) {
		return
	}

	msg := api.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = api.NewRemoveKeyboard(false)

	switch update.Message.Text {
	case helper.No:
		msg.Text = helper.MsgNoTickersDeleted

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

	case helper.Yes:
		if err = h.repo.DelSubscriptions(&user, model.EdgarSubscription); err != nil {
			log.Println(err)

			return
		}

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

		msg = api.NewMessage(update.Message.Chat.ID, helper.MsgTickersDeleted)

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
