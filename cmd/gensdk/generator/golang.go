package generator

import (
	"bytes"
	"github.com/jaronnie/jzero/cmd/gensdk/parser"
)

func init() {
	Register("go", func(language string) (Generator, error) {
		return Golang{}, nil
	})
}

type Golang struct{}

func (g Golang) Gen() ([]*GeneratedFile, error) {
	// parse
	parser.Parse()

	var files []*GeneratedFile

	files = append(files, &GeneratedFile{
		Path:    "client.go",
		Content: *bytes.NewBuffer([]byte("package client")),
	})
	return files, nil
}
