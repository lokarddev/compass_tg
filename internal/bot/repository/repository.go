package repository

import (
	"app/internal/bot/model"
	"context"
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
		{"state.previous", user.State.NavCurrent},
		{"state.current", current}}}}

	opts := options.Update().SetUpsert(true)

	coll := r.client.Database(dbName).Collection(userCollection)

	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
