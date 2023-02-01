package worker

import (
	"app/internal/bot/repository/mongo_db"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"strconv"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	// command list
	startCommand = "/start"
)

var (
	mongoURI string
)

type Worker struct {
	bot             *api.BotAPI
	botWH, apiToken string
	debug           bool
	dbClient        *mongo.Client
}

func NewWorker() (Worker, error) {
	w := Worker{}

	if err := w.initEnv(); err != nil {
		return w, err
	}

	bot, err := api.NewBotAPI(w.apiToken)
	if err != nil {
		return w, err
	}

	bot.Debug = w.debug
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

func (w *Worker) initEnv() error {
	var err error

	w.apiToken = checkEmpty(os.Getenv("API_TOKEN"))
	w.botWH = checkEmpty(os.Getenv("APP_WEBHOOK"))
	w.debug, err = strconv.ParseBool(checkEmpty(os.Getenv("DEBUG")))

	mongoURI = checkEmpty(os.Getenv("MONGODB_URI"))

	return err
}

func checkEmpty(env string) string {
	if env == "" {
		log.Fatalf("Required environment variable: %s", env)
	}

	return env
}

func (w *Worker) isValidMsg(update *api.Update) bool {
	return update.Message != nil
}

func (w *Worker) setWebhook() error {
	wh, err := api.NewWebhook(w.botWH + "/" + w.bot.Token)
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

func (w *Worker) Start() error {
	updates := w.bot.ListenForWebhook("/" + w.bot.Token)
	var err error

	go func() {
		err = http.ListenAndServe("127.0.0.1:8080", nil)
		if err != nil {
			log.Println(err)
		}
	}()

	dispatcher := NewDispatcher(w.bot)

	for update := range updates {
		if !w.isValidMsg(&update) {
			log.Println("Invalid message received")

			continue
		}

		dispatcher.CallHandler(&update)
	}

	return err
}
