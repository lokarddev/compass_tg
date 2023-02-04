package worker

import (
	"app/internal/bot/common"
	"app/internal/bot/handler"
	"app/internal/bot/handler/edgar"
	"app/internal/bot/repository"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type DispatcherInterface interface {
	AttachCommand(topic string, handler HandlersInterface)
	CallHandler(update *api.Update)
}

type HandlersInterface interface {
	Call(msg *api.Update)
}

type Dispatcher struct {
	commandHandler, textHandler map[string]HandlersInterface
}

func NewDispatcher(bot *api.BotAPI, dbClient *mongo.Client) DispatcherInterface {
	d := &Dispatcher{commandHandler: make(map[string]HandlersInterface)}

	baseRepo := repository.NewRepository(dbClient)

	d.AttachCommand(common.StartCommand, handler.NewStartHandler(bot, baseRepo))
	d.AttachCommand(common.EdgarCommand, edgar.NewEdgarHandler(bot, baseRepo))

	return d
}

func (d *Dispatcher) AttachCommand(topic string, handler HandlersInterface) {
	d.commandHandler[topic] = handler
}

func (d *Dispatcher) AttachText(text string, handler HandlersInterface) {
	d.textHandler[text] = handler
}

func (d *Dispatcher) CallHandler(update *api.Update) {
	if !d.isValidMsg(update.Message) {
		log.Println("Unavailable command used from current state")

		return
	}

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

func (d *Dispatcher) isValidMsg(msg *api.Message) bool {
	if msg.IsCommand() && msg.Text == common.StartCommand {
		return true
	}

	return true
}
