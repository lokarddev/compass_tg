package main

import (
	"app/internal/bot/worker"
	"log"
)

func main() {
	w, err := worker.NewWorker()
	if err != nil {
		log.Fatal(err)
	}

	if err = w.Start(); err != nil {
		log.Fatal(err)
	}
}
