package generator

import (
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/config"
)

func init() {
	Register("ts", func(config config.Config) (Generator, error) {
		return &Typescript{
			config: &config,
		}, nil
	})
}

type Typescript struct {
	config *config.Config
}

func (t *Typescript) Gen() ([]*GeneratedFile, error) {
	return nil, nil
}
