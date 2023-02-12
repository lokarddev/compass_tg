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
	"strings"
)

const (
	edgarSubscribeState = "EdgarSubscribe"
)

var (
	edgarSubSrc = []string{common.EdgarCommand}
)

type SubHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
	handler.BaseHandler
}

func NewSubHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *SubHandler {
	return &SubHandler{bot: bot, repo: dbRepo}
}

func (h *SubHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user.State.NavCurrent, edgarSubSrc) || update.Message.Text == keyboards.DeleteButton {
		return
	}

	inputTickers := strings.Split(strings.ReplaceAll(update.Message.Text, " ", ""), ",")

	toSub := make([]string, 0, len(inputTickers))

	for _, ticker := range inputTickers {
		cleanTicker := strings.ToUpper(ticker)
		if _, ok := handler.MockEdgarSubs[cleanTicker]; ok {
			toSub = append(toSub, cleanTicker)
		}
	}

	switch {
	case len(toSub) == 0:
		msg := api.NewMessage(update.Message.Chat.ID, message.MsgWrongTickerInput)

		if _, err = h.bot.Send(msg); err != nil {
			log.Println(err)

			return
		}

		return

	default:
		if err = h.repo.AddPendingSubs(&user, model.EdgarSubscription, toSub...); err != nil {
			log.Println(err)

			return
		}
	}

	if err = h.repo.UpsertState(&user, edgarSubscribeState); err != nil {
		log.Printf("Error setting state for user %s", user.Username)

		return
	}
}
