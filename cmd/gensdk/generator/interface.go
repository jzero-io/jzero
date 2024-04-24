package generator

import (
	"bytes"
	"github.com/pkg/errors"
)

type GeneratedFile struct {
	Path    string
	Content bytes.Buffer
}

type Generator interface {
	Gen() ([]*GeneratedFile, error)
}

var langGenerator = map[string]NewFunc{}

type NewFunc func(language string) (Generator, error)

func New(language string) (Generator, error) {
	f, ok := langGenerator[language]
	if !ok {
		return nil, errors.New("language not support")
	}
	return f(language)
}

func Register(language string, f NewFunc) {
	langGenerator[language] = f
}
