package generator

import (
	"bytes"

	"github.com/jzero-io/jzero/cmd/gensdk/config"

	"github.com/pkg/errors"
)

type GeneratedFile struct {
	Path    string
	Content bytes.Buffer
	Skip    bool
}

type Generator interface {
	Gen() ([]*GeneratedFile, error)
}

var langGenerator = map[string]NewFunc{}

type NewFunc func(target config.Config) (Generator, error)

func New(config config.Config) (Generator, error) {
	f, ok := langGenerator[config.Language]
	if !ok {
		return nil, errors.Errorf("language %s not support", config.Language)
	}
	return f(config)
}

func Register(language string, f NewFunc) {
	langGenerator[language] = f
}
