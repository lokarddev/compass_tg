package handler

import (
	"app/internal/bot/message"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DispatcherInterface interface {
	Attach(topic string, handler HandlersInterface)
	HandleUpdate(topic string, update *api.Update)
}

type HandlersInterface interface {
	Call(msg *api.Update)
}

type StartHandler struct {
	bot *api.BotAPI
}

func (h *StartHandler) Call(update *api.Update) {
	msg := api.NewMessage(update.Message.Chat.ID, message.StartReplyText)

	if _, err := h.bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func NewStartHandler(bot *api.BotAPI) *StartHandler {
	return &StartHandler{bot: bot}
}
