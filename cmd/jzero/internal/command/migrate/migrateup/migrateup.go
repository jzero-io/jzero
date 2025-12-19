package migrateup

import (
	"github.com/jzero-io/jzero/core/stores/migrate"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
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

	if len(args) > 0 {
		if cast.ToInt(args[0]) < 0 {
			return errors.New("step must be greater than 0")
		}
		return m.Up(cast.ToUint(args[0]))
	}

	return m.Up()
}
