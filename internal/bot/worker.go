package bot

import (
	"log"
	"net/http"
	"os"
	"strconv"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Worker struct {
	bot             *api.BotAPI
	botWH, apiToken string
	debug           bool
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

	return w, nil
}

func (w *Worker) initEnv() error {
	var err error

	w.apiToken = os.Getenv("API_TOKEN")
	w.botWH = os.Getenv("APP_WEBHOOK")
	w.debug, err = strconv.ParseBool(os.Getenv("DEBUG"))

	return err
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
