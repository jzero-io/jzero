package migrate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eddieowens/opts"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	DefaultXMigrationsTable = "schema_migrations"
	DefaultSource           = "file://desc/sql_migration"
)

type Driver string

const (
	MySQL  Driver = "mysql"
	Pgx    Driver = "pgx"
	Sqlite Driver = "sqlite"
)

type (
	Migrate interface {
		// Up looks at the currently active migration version
		// and will migrate all the way up (default applying all up migrations).
		Up(steps ...uint) error

		// Down looks at the currently active migration version
		// and will migrate all the way down (default applying all down migrations).
		Down(steps ...uint) error

		// Goto looks at the currently active migration version,
		// then migrates either up or down to the specified version.
		Goto(version uint) error

		// Version returns the currently active migration version.
		// If no migration has been applied, yet, it will return ErrNilVersion.
		Version() (version uint, dirty bool, err error)

		// Close source and database, return source error and database error
		Close() (error, error)
	}

	MigrateOpts struct {
		Source             string
		SourceAppendDriver bool
		XMigrationsTable   string
	}

	defaultMigrate struct {
		migrate *migrate.Migrate
	}
)

func (d *defaultMigrate) Close() (error, error) {
	return d.migrate.Close()
}

func WithSource(source string) opts.Opt[MigrateOpts] {
	return func(d *MigrateOpts) {
		d.Source = source
	}
}

func WithSourceAppendDriver(sourceAppendDriver bool) opts.Opt[MigrateOpts] {
	return func(d *MigrateOpts) {
		d.SourceAppendDriver = sourceAppendDriver
	}
}

func WithXMigrationsTable(xMigrationsTable string) opts.Opt[MigrateOpts] {
	return func(u *MigrateOpts) {
		u.XMigrationsTable = xMigrationsTable
	}
}

func (d MigrateOpts) DefaultOptions() MigrateOpts {
	return MigrateOpts{
		Source:             DefaultSource,
		XMigrationsTable:   DefaultXMigrationsTable,
		SourceAppendDriver: false,
	}
}

func NewMigrate(sqlConf sqlx.SqlConf, op ...opts.Opt[MigrateOpts]) (Migrate, error) {
	ops := opts.DefaultApply(op...)

	var (
		dataSource     = sqlConf.DataSource
		source         = ops.Source
		paramConnector string
	)

	if strings.Contains(dataSource, "?") {
		paramConnector = "&"
	} else {
		paramConnector = "?"
	}

	switch Driver(sqlConf.DriverName) {
	case MySQL:
		dataSource = "mysql://" + dataSource
	case Pgx:
		dataSource = "pgx5://" + strings.TrimPrefix(dataSource, "postgres://")
	case Sqlite:
		dataSource = "sqlite://" + dataSource
	default:
		return nil, fmt.Errorf("unsupported driver: %s", sqlConf.DriverName)
	}

	if ops.SourceAppendDriver {
		source = fmt.Sprintf("%s/%s", source, sqlConf.DriverName)
	}

	dataSource = fmt.Sprintf("%s%sx-migrations-table=%s", dataSource, paramConnector, ops.XMigrationsTable)

	m, err := migrate.New(source, dataSource)
	if err != nil {
		return nil, err
	}

	return &defaultMigrate{
		migrate: m,
	}, nil
}

func (d *defaultMigrate) Up(steps ...uint) error {
	if len(steps) > 1 {
		return errors.New("steps number should not be more than 1")
	}

	var err error

	if len(steps) == 0 {
		err = d.migrate.Up()
	} else {
		err = d.migrate.Steps(int(steps[0]))
	}

	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}

func (d *defaultMigrate) Down(steps ...uint) error {
	if len(steps) > 1 {
		return errors.New("steps number should not be more than 1")
	}

	if len(steps) == 0 {
		return d.migrate.Down()
	}

	return d.migrate.Steps(-int(steps[0]))
}

func (d *defaultMigrate) Goto(version uint) error {
	return d.migrate.Migrate(version)
}

func (d *defaultMigrate) Version() (uint, bool, error) {
	return d.migrate.Version()
}
