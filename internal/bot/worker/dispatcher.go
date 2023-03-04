package worker

import (
	"app/internal/bot/handler"
	"app/internal/bot/handler/edgar"
	"app/internal/bot/helper"
	"app/internal/repository"
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DispatcherInterface interface {
	AttachCommand(topic string, handler HandlersInterface)
	AttachText(handler HandlersInterface)
	CallHandler(update *api.Update)
}

type HandlersInterface interface {
	Call(msg *api.Update)
}

type Dispatcher struct {
	commandHandler map[string]HandlersInterface
	textHandler    []HandlersInterface
}

func NewDispatcher(bot *api.BotAPI, repo repository.BaseRepoInterface) DispatcherInterface {
	d := &Dispatcher{commandHandler: make(map[string]HandlersInterface)}

	// helper commands
	d.AttachCommand(helper.StartCommand, handler.NewStartHandler(bot, repo))

	// edgar text/command handlers
	d.AttachCommand(helper.EdgarCommand, edgar.NewSubCheckHandler(bot, repo))
	d.AttachText(edgar.NewSubHandler(bot, repo, helper.EdgarBranch))
	d.AttachText(edgar.NewSubApproveHandler(bot, repo, helper.EdgarBranch))
	d.AttachText(edgar.NewSubFinalHandler(bot, repo, helper.EdgarBranch))
	d.AttachText(edgar.NewDeleteHandler(bot, repo, helper.EdgarBranch))
	d.AttachText(edgar.NewDeleteChoiceHandler(bot, repo, helper.EdgarBranch))
	d.AttachText(edgar.NewDelSingleApproveHandler(bot, repo, helper.EdgarBranch))
	d.AttachText(edgar.NewDelAllApproveHandler(bot, repo, helper.EdgarBranch))
	d.AttachText(edgar.NewDelSingleFinalHandler(bot, repo, helper.EdgarBranch))
	d.AttachText(edgar.NewDelAllFinalHandler(bot, repo, helper.EdgarBranch))

	return d
}

func (d *Dispatcher) AttachCommand(topic string, handler HandlersInterface) {
	d.commandHandler[topic] = handler
}

func (d *Dispatcher) AttachText(handler HandlersInterface) {
	d.textHandler = append(d.textHandler, handler)
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
		for _, h := range d.textHandler {
			h.Call(update)
		}
	}
}

func (d *Dispatcher) isValidMsg(msg *api.Message) bool {
	if msg.IsCommand() && msg.Text == helper.StartCommand {
		return true
	}

	return true
}
