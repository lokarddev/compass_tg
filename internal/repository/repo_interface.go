package repository

import (
	"app/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BaseRepoInterface interface {
	UpsertUser(sentFrom *tgbotapi.User) (model.User, error)
	UpsertState(user *model.User, current, branch int) error
	GetUser(sentFrom *tgbotapi.User) (model.User, error)
	AddPendingSubs(user *model.User, subType string, subs ...string) error
	UpsertSubscriptions(user *model.User, subType string, subs ...string) error
	DelSubscriptions(user *model.User, subType string, subs ...string) error
}
