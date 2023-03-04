package main

import (
	"app/internal/bot/worker"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	w, err := worker.NewWorker()
	if err != nil {
		log.Fatal(err)
	}

	if err = w.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
