package main

import (
	"app/internal/bot"
	"log"
)

func main() {
	worker, err := bot.NewWorker()
	if err != nil {
		log.Fatal(err)
	}

	if err = worker.Start(); err != nil {
		log.Fatal(err)
	}
}
