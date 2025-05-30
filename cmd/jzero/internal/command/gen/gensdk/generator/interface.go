package generator

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/config"
	gconfig "github.com/jzero-io/jzero/cmd/jzero/internal/config"
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
	f, ok := langGenerator[gconfig.C.Gen.Sdk.Language]
	if !ok {
		return nil, errors.Errorf("language %s not support", gconfig.C.Gen.Sdk.Language)
	}
	return f(config)
}

func Register(language string, f NewFunc) {
	langGenerator[language] = f
}
