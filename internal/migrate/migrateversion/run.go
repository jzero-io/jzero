package migrateversion

import (
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jzero-io/jzero/config"
)

func Run(args []string) error {
	m, err := migrate.New(config.C.Migrate.Source, config.C.Migrate.Database)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil {
		return err
	}

	if dirty {
		fmt.Printf("%v (dirty)\n", version)
	} else {
		fmt.Printf("%v\n", version)
	}
	return nil
}
