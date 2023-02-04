package worker

import (
	"app/internal/bot/repository/mongo_db"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
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

func (w *Worker) Start() error {
	updates := w.bot.ListenForWebhook("/" + w.bot.Token)
	var err error

	go func() {
		err = http.ListenAndServe("127.0.0.1:8080", nil)
		if err != nil {
			log.Println(err)
		}
	}()

	dispatcher := NewDispatcher(w.bot, w.dbClient)

	for update := range updates {
		if w.isEmptyMsg(&update) {
			log.Println("Empty/nil message received")

			continue
		}

		dispatcher.CallHandler(&update)
	}

	return err
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
