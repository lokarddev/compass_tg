package edgar

import (
	"app/internal/bot/common"
	"app/internal/bot/common/keyboards"
	"app/internal/bot/common/message"
	"app/internal/bot/handler"
	"app/internal/bot/model"
	"app/internal/bot/repository/mongo_db"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var (
	edgarDelAllFinalSrc = []string{edgarDelAllApproveState}
)

type DelAllFinalHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
	handler.BaseHandler
}

func NewDelAllFinalHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *DelAllFinalHandler {
	return &DelAllFinalHandler{bot: bot, repo: dbRepo}
}

func (h *DelAllFinalHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user.State.NavCurrent, edgarDelAllFinalSrc) || !h.ValidText(update.Message.Text, edgarSubFinalButtons) {
		return
	}

	msg := api.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = api.NewRemoveKeyboard(false)

	switch update.Message.Text {
	case keyboards.No:
		msg.Text = message.MsgNoTickersDeleted

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

	case keyboards.Yes:
		if err = h.repo.DelSubscriptions(&user, model.EdgarSubscription); err != nil {
			log.Println(err)

			return
		}

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

		msg = api.NewMessage(update.Message.Chat.ID, message.MsgTickersDeleted)

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
