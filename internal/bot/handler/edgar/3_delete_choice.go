package edgar

import (
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
	edgarDeleteSingleState = "edgarDeleteSingle"
	edgarDeleteAllState    = "edgarDeleteAll"
)

var (
	edgarDelSrc           = []string{edgarDeleteState}
	edgarDeleteAllButtons = []string{keyboards.DeleteAllButton}
)

type DeleteChoiceHandler struct {
	bot  *api.BotAPI
	repo mongo_db.BaseRepoInterface
	handler.BaseHandler
}

func NewDeleteChoiceHandler(bot *api.BotAPI, dbRepo mongo_db.BaseRepoInterface) *DeleteChoiceHandler {
	return &DeleteChoiceHandler{bot: bot, repo: dbRepo}
}

func (h *DeleteChoiceHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user.State.NavCurrent, edgarDelSrc) || update.Message.Text == keyboards.DeleteButton {
		return
	}

	switch {
	case h.ValidText(update.Message.Text, edgarDeleteAllButtons):
		if err = h.repo.AddPendingSubs(&user, model.EdgarSubscription, user.Subscriptions.Edgar.Tickers...); err != nil {
			log.Println(err)

			return
		}

		if err = h.repo.UpsertState(&user, edgarDeleteAllState); err != nil {
			log.Printf("Error setting state for user %s", user.Username)

			return
		}

	default:
		inputTickers := strings.Split(strings.ReplaceAll(update.Message.Text, " ", ""), ",")

		toDel := make([]string, 0, len(inputTickers))

		for _, ticker := range inputTickers {
			cleanTicker := strings.ToUpper(ticker)
			if _, ok := handler.MockEdgarSubs[cleanTicker]; ok {
				toDel = append(toDel, cleanTicker)
			}
		}

		switch {
		case len(toDel) == 0:
			msg := api.NewMessage(update.Message.Chat.ID, message.MsgWrongTickerInput)

			if _, err = h.bot.Send(msg); err != nil {
				log.Println(err)

				return
			}

			return

		default:
			if err = h.repo.AddPendingSubs(&user, model.EdgarSubscription, toDel...); err != nil {
				log.Println(err)

				return
			}
		}

		if err = h.repo.UpsertState(&user, edgarDeleteSingleState); err != nil {
			log.Printf("Error setting state for user %s", user.Username)

			return
		}
	}
}
