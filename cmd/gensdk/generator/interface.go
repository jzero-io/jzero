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

type Target struct {
	Language string
	APP      string
	Module   string
}

type NewFunc func(target Target) (Generator, error)

func New(target Target) (Generator, error) {
	f, ok := langGenerator[target.Language]
	if !ok {
		return nil, errors.Errorf("language %s not support", target.Language)
	}
	return f(target)
}

func Register(language string, f NewFunc) {
	langGenerator[language] = f
}
