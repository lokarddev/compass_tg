package edgar

import (
	"app/internal/bot/common"
	"app/internal/bot/common/keyboards"
	"app/internal/bot/common/message"
	"app/internal/bot/handler"
	"app/internal/bot/model"
	"app/internal/bot/repository/mongo_db"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	_ = "EdgarSubscribeFinal"
)

var (
	edgarSubFinalSrc     = []string{edgarSubscribeApproveState}
	edgarSubFinalButtons = []string{keyboards.Yes, keyboards.No}
)

type SubFinalHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
	handler.BaseHandler
}

func NewSubFinalHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *SubFinalHandler {
	return &SubFinalHandler{bot: bot, repo: dbRepo}
}

func (h *SubFinalHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user.State.NavCurrent, edgarSubFinalSrc) || !h.ValidText(update.Message.Text, edgarSubFinalButtons) {
		return
	}

	msg := api.NewMessage(update.Message.Chat.ID, message.MsgSubSuccess)
	msg.ReplyMarkup = api.NewRemoveKeyboard(false)

	switch update.Message.Text {
	case keyboards.No:
		msg.Text = message.MsgNoSub

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

	case keyboards.Yes:
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

	if err = h.repo.UpsertState(&user, common.StartCommand); err != nil {
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
