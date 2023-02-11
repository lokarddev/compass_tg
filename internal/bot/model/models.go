package model

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	EdgarSubscription = "edgar"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	TgID          int64              `json:"tg_id" bson:"tg_id"`
	FirstName     string             `json:"first_name" bson:"first_name"`
	LastName      string             `json:"last_name" bson:"last_name"`
	Username      string             `json:"username" bson:"username"`
	LangCode      string             `json:"lang_code" bson:"lang_code"`
	State         State              `json:"state" bson:"state"`
	Subscriptions Subscriptions      `json:"subscriptions" bson:"subscriptions"`
}

type State struct {
	NavPrevious string `json:"previous" bson:"previous"`
	NavCurrent  string `json:"current" bson:"current"`
}

type Subscriptions struct {
	Edgar          Subscription `json:"edgar" bson:"edgar"`
	CompanyData    Subscription `json:"company_data" bson:"company_data"`
	InsiderTrading Subscription `json:"insider_trading" bson:"insider_trading"`
	PriceAlert     Subscription `json:"price_alert" bson:"price_alert"`
}

type Subscription struct {
	Enabled     bool     `json:"enabled" bson:"enabled"`
	PendingSubs []string `json:"pending_subs" bson:"pending_subs"`
	Tickers     []string `json:"tickers" bson:"tickers"`
}
