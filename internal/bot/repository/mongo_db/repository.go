package mongo_db

import (
	"app/internal/bot/model"
	"app/internal/bot/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "compass_tg"

	userCollection = "users"
)

type BaseRepoInterface interface {
	UpsertUser(sentFrom *tgbotapi.User) (model.User, error)
	UpsertState(user *model.User, current string) error
	GetUser(sentFrom *tgbotapi.User) (model.User, error)
	AddPendingSubs(user *model.User, subType string, subs ...string) error
	UpsertSubscriptions(user *model.User, subType string, subs ...string) error
	DelSubscriptions(user *model.User, subType string, subs ...string) error
}

type BaseRepository struct {
	client *mongo.Client
}

func NewRepository(dbClient *mongo.Client) *BaseRepository {
	return &BaseRepository{client: dbClient}
}

func (r *BaseRepository) GetUser(sentFrom *tgbotapi.User) (model.User, error) {
	var user model.User

	filter := bson.D{{"tg_id", sentFrom.ID}}
	coll := r.client.Database(dbName).Collection(userCollection)

	if err := coll.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *BaseRepository) UpsertUser(sentFrom *tgbotapi.User) (model.User, error) {
	var user model.User

	filter := bson.D{{"tg_id", sentFrom.ID}}
	update := bson.D{{"$set", bson.D{
		{"first_name", sentFrom.FirstName},
		{"last_name", sentFrom.LastName},
		{"username", sentFrom.UserName},
		{"lang_code", sentFrom.LanguageCode}}}}

	opts := options.Update().SetUpsert(true)

	coll := r.client.Database(dbName).Collection(userCollection)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return user, err
	}

	filter = bson.D{{"tg_id", sentFrom.ID}}

	if err = coll.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (r *BaseRepository) UpsertState(user *model.User, current string) error {
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{"state.current", current}}}}

	opts := options.Update().SetUpsert(true)

	coll := r.client.Database(dbName).Collection(userCollection)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepository) AddPendingSubs(user *model.User, subType string, subs ...string) error {
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{fmt.Sprintf("subscriptions.%s.pending_subs", subType), subs}}}}

	opts := options.Update().SetUpsert(true)

	coll := r.client.Database(dbName).Collection(userCollection)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepository) UpsertSubscriptions(user *model.User, subType string, subs ...string) error {
	switch subType {
	case model.EdgarSubscription:
		newEdgarSubs(user, subs)
	}

	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{fmt.Sprintf("subscriptions.%s.tickers", subType), user.Subscriptions.Edgar.Tickers},
		{fmt.Sprintf("subscriptions.%s.enabled", subType), user.Subscriptions.Edgar.Enabled},
		{fmt.Sprintf("subscriptions.%s.pending_subs", subType), user.Subscriptions.Edgar.PendingSubs}}}}

	opts := options.Update().SetUpsert(true)

	coll := r.client.Database(dbName).Collection(userCollection)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func newEdgarSubs(user *model.User, subs []string) {
	switch {
	case len(subs) > 0:
		user.Subscriptions.Edgar.Tickers = utils.RemoveDuplicatesAndSortStr(append(user.Subscriptions.Edgar.Tickers, subs...))
	default:
		user.Subscriptions.Edgar.Tickers = utils.RemoveDuplicatesAndSortStr(append(user.Subscriptions.Edgar.Tickers, user.Subscriptions.Edgar.PendingSubs...))
	}

	if len(user.Subscriptions.Edgar.Tickers) > 0 {
		user.Subscriptions.Edgar.Enabled = true
	}

	user.Subscriptions.Edgar.PendingSubs = []string{}
}

func (r *BaseRepository) DelSubscriptions(user *model.User, subType string, subs ...string) error {
	switch subType {
	case model.EdgarSubscription:
		edgarSetSubs(user, subs)
	}

	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{fmt.Sprintf("subscriptions.%s.tickers", subType), user.Subscriptions.Edgar.Tickers},
		{fmt.Sprintf("subscriptions.%s.enabled", subType), user.Subscriptions.Edgar.Enabled},
		{fmt.Sprintf("subscriptions.%s.pending_subs", subType), user.Subscriptions.Edgar.PendingSubs}}}}

	opts := options.Update().SetUpsert(true)

	coll := r.client.Database(dbName).Collection(userCollection)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func edgarSetSubs(user *model.User, subs []string) {
	toSet := make([]string, 0, len(user.Subscriptions.Edgar.Tickers))

	for _, ticker := range subs {
		if !utils.StrInSlice(user.Subscriptions.Edgar.Tickers, ticker) {
			toSet = append(toSet, ticker)
		}
	}

	switch {
	case len(toSet) > 0:
		user.Subscriptions.Edgar.Tickers = utils.RemoveDuplicatesAndSortStr(append(user.Subscriptions.Edgar.Tickers, toSet...))
	default:
		user.Subscriptions.Edgar.Tickers = []string{}
	}

	if len(toSet) > 0 {
		user.Subscriptions.Edgar.Enabled = true
	} else {
		user.Subscriptions.Edgar.Enabled = false
	}

	user.Subscriptions.Edgar.PendingSubs = []string{}
}
