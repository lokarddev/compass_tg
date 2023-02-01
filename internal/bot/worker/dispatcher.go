package worker

import (
	"app/internal/bot/handler"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Dispatcher struct {
	commandHandler, textHandler map[string]handler.HandlerInterface
}

func (d *Dispatcher) AttachCommand(topic string, handler handler.HandlerInterface) {
	d.commandHandler[topic] = handler
}

func (d *Dispatcher) AttachText(text string, handler handler.HandlerInterface) {
	d.textHandler[text] = handler
}

func (d *Dispatcher) CallHandler(update *api.Update) {
	switch update.Message.IsCommand() {
	case true:
		h, ok := d.commandHandler[update.Message.Text]
		if ok {
			h.Call(update)
		}

	case false:
		h, ok := d.textHandler[update.Message.Text]
		if ok {
			h.Call(update)
		}
	}
}

func NewDispatcher(bot *api.BotAPI) *Dispatcher {
	d := &Dispatcher{commandHandler: make(map[string]handler.HandlerInterface)}

	d.AttachCommand(startCommand, handler.NewStartHandler(bot))

	return d
}
