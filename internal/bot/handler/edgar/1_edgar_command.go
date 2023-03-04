package edgar

import (
	"app/internal/bot/handler"
	"app/internal/bot/helper"
	"app/internal/repository"
	"fmt"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SubCheckHandler struct {
	bot  *api.BotAPI
	repo repository.BaseRepoInterface
	handler.BaseHandler
}

func NewSubCheckHandler(bot *api.BotAPI, dbRepo repository.BaseRepoInterface) *SubCheckHandler {
	return &SubCheckHandler{bot: bot, repo: dbRepo}
}

func (h *SubCheckHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	msg := api.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = api.NewRemoveKeyboard(false)

	switch user.Subscriptions.Edgar.Enabled {

	case true:
		var subs string
		for _, sub := range user.Subscriptions.Edgar.Tickers {
			subs += fmt.Sprintf("%s\n", sub)
		}

		msg.Text = fmt.Sprintf("%s\n%s", helper.MsgEdgarSubs, subs)
		msg.ReplyMarkup = helper.DeleteSubKeyboard

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

	case false:
		msg.Text = helper.MsgHasNoSubscriptions

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}
	}

	if _, err = h.bot.Send(api.NewMessage(update.Message.Chat.ID, helper.MsgTypeTickersToSubscribe)); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, helper.EdgarCommandState, helper.EdgarBranch); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
