package edgar

import (
	"app/internal/bot/common"
	"app/internal/bot/common/keyboards"
	"app/internal/bot/common/message"
	"app/internal/bot/handler"
	"app/internal/bot/repository/mongo_db"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type SubCheckHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
	handler.BaseHandler
}

func NewSubCheckHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *SubCheckHandler {
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

		msg.Text = fmt.Sprintf("%s\n%s", message.MsgEdgarSubs, subs)
		msg.ReplyMarkup = keyboards.DeleteSubKeyboard

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

	case false:
		msg.Text = message.MsgHasNoSubscriptions

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}
	}

	if _, err = h.bot.Send(api.NewMessage(update.Message.Chat.ID, message.MsgTypeTickersToSubscribe)); err != nil {
		log.Println(err)

		return
	}

	if err = h.repo.UpsertState(&user, common.EdgarCommand); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
