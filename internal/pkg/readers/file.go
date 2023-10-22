package readers

import (
	"bufio"
	"os"
	"strings"
)

type FileService interface {
	ReadFile(fileLoc string) ([]string, error)
}

type file struct{}

var _ FileService = (*file)(nil)

func NewFile() *file {
	return &file{}
}

func (f *file) ReadFile(fileLoc string) ([]string, error) {

	file, err := os.Open(fileLoc)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.TrimSpace(line)
		line = strings.ToLower(line)
		fileLines = append(fileLines, line)
	}

	file.Close()

	return fileLines, nil
}
