package main

import (
	"app/internal/bot/worker"
	"log"
)

func main() {
	worker, err := worker.NewWorker()
	if err != nil {
		log.Fatal(err)
	}

	if err = worker.Start(); err != nil {
		log.Fatal(err)
	}
}
