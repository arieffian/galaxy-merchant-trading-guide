package main

import (
	"context"
	"time"

	"github.com/arieffian/roman-alien-currency/internal/app"
	"github.com/arieffian/roman-alien-currency/internal/pkg/converters"
	"github.com/arieffian/roman-alien-currency/internal/pkg/parsers"
	"github.com/arieffian/roman-alien-currency/internal/pkg/readers"
	log "github.com/sirupsen/logrus"
)

const (
	contextDeadline = 10 * time.Second
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	ctx, cancel := context.WithTimeout(context.Background(), contextDeadline)
	defer cancel()

	converter := converters.NewConverter()
	parser := parsers.NewParser(parsers.NewParserParams{
		Converter: converter,
	})
	fileReader := readers.NewFile()

	cli, err := app.NewCli(app.NewCliParams{
		Converter:  converter,
		Parser:     parser,
		FileReader: fileReader,
	})

	if err != nil {
		log.Fatalf("failed to create the new cli: %s\n", err)
	}

	err = cli.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
