package main

import (
	"context"
	"time"

	"github.com/arieffian/roman-alien-currency/internal/app"
	log "github.com/sirupsen/logrus"
)

const (
	contextDeadline = 10 * time.Second
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	ctx, cancel := context.WithTimeout(context.Background(), contextDeadline)
	defer cancel()

	cli, err := app.NewCli(ctx)
	if err != nil {
		log.Fatalf("failed to create the new cli: %s\n", err)
	}

	err = cli.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
