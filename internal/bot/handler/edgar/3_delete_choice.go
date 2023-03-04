package edgar

import (
	"app/internal/bot/handler"
	"app/internal/bot/helper"
	"app/internal/model"
	"app/internal/repository"
	"log"
	"strings"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var edgarDeleteAllButtons = []string{helper.DeleteAllButton}

type DeleteChoiceHandler struct {
	bot  *api.BotAPI
	repo repository.BaseRepoInterface
	handler.BaseHandler
	nav handler.Nav
}

func NewDeleteChoiceHandler(bot *api.BotAPI, dbRepo repository.BaseRepoInterface, branch int) *DeleteChoiceHandler {
	return &DeleteChoiceHandler{
		bot:  bot,
		repo: dbRepo,
		nav: handler.Nav{
			ValidSources: []int{helper.EdgarDeleteState},
			ValidBranch:  branch,
		}}
}

func (h *DeleteChoiceHandler) Call(update *api.Update) {
	user, err := h.repo.GetUser(update.SentFrom())
	if err != nil {
		log.Printf("Cannot get user, due to error: %s", err.Error())

		return
	}

	if !h.ValidState(user, h.nav) || update.Message.Text == helper.DeleteButton {
		return
	}

	switch {
	case h.ValidText(update.Message.Text, edgarDeleteAllButtons):
		if err = h.repo.AddPendingSubs(&user, model.EdgarSubscription, user.Subscriptions.Edgar.Tickers...); err != nil {
			log.Println(err)

			return
		}

		if err = h.repo.UpsertState(&user, helper.EdgarDeleteAllState, helper.EdgarBranch); err != nil {
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
			msg := api.NewMessage(update.Message.Chat.ID, helper.MsgWrongTickerInput)

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

		if err = h.repo.UpsertState(&user, helper.EdgarDeleteSingleState, helper.EdgarBranch); err != nil {
			log.Printf("Error setting state for user %s", user.Username)

			return
		}
	}
}
