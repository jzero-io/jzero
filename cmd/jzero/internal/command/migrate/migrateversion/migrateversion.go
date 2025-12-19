package migrateversion

import (
	"fmt"

	"github.com/jzero-io/jzero/core/stores/migrate"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

func Run(args []string) error {
	m, err := migrate.NewMigrate(sqlx.SqlConf{
		DataSource: config.C.Migrate.DataSourceUrl,
		DriverName: config.C.Migrate.Driver,
	},
		migrate.WithXMigrationsTable(config.C.Migrate.XMigrationsTable),
		migrate.WithSource(config.C.Migrate.Source),
		migrate.WithSourceAppendDriver(config.C.Migrate.SourceAppendDriver))
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
