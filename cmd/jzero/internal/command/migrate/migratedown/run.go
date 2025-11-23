package migratedown

import (
	"errors"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cast"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

func Run(args []string) error {
	m, err := migrate.New(config.C.Migrate.Source, config.C.Migrate.Database)
	if err != nil {
		return err
	}

	if len(args) > 0 {
		if cast.ToInt(args[0]) < 0 {
			return errors.New("step must be greater than 0")
		}
		return m.Steps(-cast.ToInt(args[0]))
	}
	return m.Steps(-1)
}
