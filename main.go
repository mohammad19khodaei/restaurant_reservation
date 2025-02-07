package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohammad19khodaei/restaurant_reservation/config"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/application"
)

func main() {
	config, err := config.LoadConfig("./config", "config")
	if err != nil {
		log.Fatal("could not load config: ", err)
	}

	app, err := application.New(config)
	if err != nil {
		log.Fatal("could not create application: ", err)
	}

	app.RegisterRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.Run(ctx)

	closeSignal := make(chan os.Signal, 1)
	signal.Notify(closeSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-closeSignal:
	case <-ctx.Done():
	}
}
