package worker

import (
	"log"
	"os"
	"strconv"
)

var (
	apiToken string
	botWH    string
	debug    bool
	mongoURI string
)

func (w *Worker) initEnv() error {
	var err error

	apiToken = checkEmpty(os.Getenv("API_TOKEN"))
	botWH = checkEmpty(os.Getenv("APP_WEBHOOK"))
	debug, err = strconv.ParseBool(checkEmpty(os.Getenv("DEBUG")))
	mongoURI = checkEmpty(os.Getenv("MONGODB_URI"))

	return err
}

func checkEmpty(env string) string {
	if env == "" {
		log.Fatalf("Required environment variable: %s", env)
	}

	return env
}
