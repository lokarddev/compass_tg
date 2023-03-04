package worker

import (
	"app/internal/repository/mongo_db"
	"context"
	"log"
	"net/http"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go.mongodb.org/mongo-driver/mongo"
)

type Worker struct {
	bot      *api.BotAPI
	dbClient *mongo.Client
}

func NewWorker() (Worker, error) {
	w := Worker{}

	if err := w.initEnv(); err != nil {
		return w, err
	}

	bot, err := api.NewBotAPI(apiToken)
	if err != nil {
		return w, err
	}

	bot.Debug = debug
	w.bot = bot

	log.Printf("Authorized on account %s", bot.Self.UserName)

	if err = w.setWebhook(); err != nil {
		return w, err
	}

	w.dbClient, err = mongo_db.NewClient(mongoURI)
	if err != nil {
		return w, err
	}

	return w, nil
}

func (w *Worker) Start(ctx context.Context) error {
	baseRepo := mongo_db.NewRepository(ctx, w.dbClient)

	dispatcher := NewDispatcher(w.bot, baseRepo)

	return http.ListenAndServe("127.0.0.1:8080", w.WebhookHandler(dispatcher))
}

func (w *Worker) isEmptyMsg(update *api.Update) bool {
	return update.Message == nil
}

func (w *Worker) setWebhook() error {
	wh, err := api.NewWebhook(botWH + "/" + w.bot.Token)
	if err != nil {
		log.Println(err)

		return err
	}

	_, err = w.bot.Request(wh)
	if err != nil {
		log.Println(err)

		return err
	}

	info, err := w.bot.GetWebhookInfo()
	if err != nil {
		log.Println(err)

		return err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return err
}

func (w *Worker) WebhookHandler(dispatcher DispatcherInterface) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		update, err := w.bot.HandleUpdate(request)
		if err != nil {
			log.Println(err)
			return
		}

		if w.isEmptyMsg(update) {
			log.Println("Empty/nil message received")

			return
		}

		dispatcher.CallHandler(update)
	}
}
