package handler

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DispatcherInterface interface {
	Attach(topic string, handler HandlerInterface)
	HandleUpdate(topic string, update *api.Update)
}

type HandlerInterface interface {
	Call(msg *api.Update)
}
