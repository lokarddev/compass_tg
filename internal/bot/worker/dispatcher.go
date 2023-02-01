package worker

import (
	"app/internal/bot/handlers"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Dispatcher struct {
	commandHandler, textHandler map[string]handlers.HandlerInterface
}

func (d *Dispatcher) AttachCommand(topic string, handler handlers.HandlerInterface) {
	d.commandHandler[topic] = handler
}

func (d *Dispatcher) AttachText(text string, handler handlers.HandlerInterface) {
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
	d := &Dispatcher{commandHandler: make(map[string]handlers.HandlerInterface)}

	d.AttachCommand(startCommand, handlers.NewStartHandler(bot))

	return d
}
