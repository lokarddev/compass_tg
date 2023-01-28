package bot

import (
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	// command list
	startCommand = "/start"
)

type DispatcherInterface interface {
	Attach(topic string, handler HandlerInterface)
	HandleUpdate(topic string, update *api.Update)
}

type HandlerInterface interface {
	Call(msg *api.Update)
}

type Dispatcher struct {
	commandHandler, textHandler map[string]HandlerInterface
}

func (d *Dispatcher) AttachCommand(topic string, handler HandlerInterface) {
	d.commandHandler[topic] = handler
}

func (d *Dispatcher) AttachText(text string, handler HandlerInterface) {
	d.textHandler[text] = handler
}

func (d *Dispatcher) CallHandler(update *api.Update) {
	switch update.Message.IsCommand() {
	case true:
		handler, ok := d.commandHandler[update.Message.Text]
		if ok {
			handler.Call(update)
		}

	case false:
		handler, ok := d.textHandler[update.Message.Text]
		if ok {
			handler.Call(update)
		}
	}
}

func NewDispatcher(bot *api.BotAPI) *Dispatcher {
	d := &Dispatcher{commandHandler: make(map[string]HandlerInterface)}

	d.AttachCommand(startCommand, NewStartHandler(bot))

	return d
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
