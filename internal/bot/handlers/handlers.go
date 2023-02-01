package handlers

import (
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DispatcherInterface interface {
	Attach(topic string, handler HandlerInterface)
	HandleUpdate(topic string, update *api.Update)
}

type HandlerInterface interface {
	Call(msg *api.Update)
}

type StartHandler struct {
	bot *api.BotAPI
}

func (h *StartHandler) Call(update *api.Update) {
	msg := api.NewMessage(update.Message.Chat.ID, startReplyText)

	if _, err := h.bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func NewStartHandler(bot *api.BotAPI) *StartHandler {
	return &StartHandler{bot: bot}
}
