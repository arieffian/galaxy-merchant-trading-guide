package app

import (
	"context"
	"fmt"
	"strings"

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

	lines, err := c.fileReader.ReadFile("input")
	if err != nil {
		return err
	}

	indices := []int{}
	for idx, line := range lines {
		lineArr := strings.Split(line, " ")

		found := c.parser.ParseCurrency(lineArr)
		if found {
			indices = append(indices, idx)
		}
	}

	// remove currency from params
	for i, idx := range indices {
		lines = append(lines[:idx-i], lines[idx+1-i:]...)
	}

	indices = []int{}
	for idx, line := range lines {
		lineArr := strings.Split(line, " ")

		found, err := c.parser.ParseMetal(lineArr)
		if err != nil {
			return err
		}
		if found {
			indices = append(indices, idx)
		}
	}

	// remove currency from params
	for i, idx := range indices {
		lines = append(lines[:idx-i], lines[idx+1-i:]...)
	}

	answers, _ := c.parser.ProcessQuestion(lines)

	for _, answer := range answers {
		fmt.Println(answer)
	}

	return nil
}
