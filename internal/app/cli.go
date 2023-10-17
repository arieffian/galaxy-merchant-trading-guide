package app

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arieffian/roman-alien-currency/internal/pkg/converters"
)

type cli struct {
	converter converters.ConverterService
}

func NewCli(ctx context.Context) (*cli, error) {
	converter := converters.NewConverter()

	return &cli{
		converter: converter,
	}, nil
}

func (c *cli) Run(ctx context.Context) error {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter file location: ")
	fileLoc, _ := reader.ReadString('\n')

	fileLoc = fileLoc[:len(fileLoc)-1]

	file, err := os.Open(fileLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println(c.converter.ArabicToRoman(999))
	fmt.Println(c.converter.RomanToArabic("XXX"))

	return nil
}
