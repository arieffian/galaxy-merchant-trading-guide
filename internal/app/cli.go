package app

import (
	"context"
	"fmt"

	"github.com/arieffian/roman-alien-currency/internal/pkg/converters"
	"github.com/arieffian/roman-alien-currency/internal/pkg/parsers"
	"github.com/arieffian/roman-alien-currency/internal/pkg/readers"
)

type cli struct {
	converter  converters.ConverterService
	parser     parsers.ParserService
	fileReader readers.FileService
}

type NewCliParams struct {
	Converter  converters.ConverterService
	Parser     parsers.ParserService
	FileReader readers.FileService
}

func NewCli(p NewCliParams) (*cli, error) {

	return &cli{
		converter:  p.Converter,
		parser:     p.Parser,
		fileReader: p.FileReader,
	}, nil
}

func (c *cli) Run(ctx context.Context) error {

	body, err := c.fileReader.ReadFile("input")
	if err != nil {
		return err
	}

	results, err := c.parser.Parse(body)
	if err != nil {
		return err
	}

	for _, result := range results {
		fmt.Println(result)
	}

	return nil
}
