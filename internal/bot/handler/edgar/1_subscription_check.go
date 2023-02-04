package edgar

import (
	"app/internal/bot/common"
	"app/internal/bot/common/message"
	"app/internal/bot/handler"
	"app/internal/bot/repository"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type SubCheckHandler struct {
	bot  *api.BotAPI
	repo repository.BaseRepoInterface
	handler.BaseHandler
}

func NewEdgarHandler(bot *api.BotAPI, dbRepo repository.BaseRepoInterface) *SubCheckHandler {
	return &SubCheckHandler{bot: bot, repo: dbRepo}
}

func (h *SubCheckHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if err = h.repo.UpsertState(&user, common.StartCommand); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}

	switch user.Subscriptions.Edgar.Enabled {

	case true:
		var subs string
		for _, sub := range user.Subscriptions.Edgar.Tickers {
			subs += fmt.Sprintf("%s\n", sub)
		}

		msg := api.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n %s", message.MsgEdgarSubs, subs))
		msg.ReplyMarkup = common.DeleteSubKeyboard

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)
		}

	case false:
		if _, err = h.bot.Send(api.NewMessage(update.Message.Chat.ID, message.MsgHasNoSubscriptions)); err != nil {
			log.Println(err)
		}
	}

	if _, err = h.bot.Send(api.NewMessage(update.Message.Chat.ID, message.MsgTypeTickersToSubscribe)); err != nil {
		log.Println(err)
	}
}
