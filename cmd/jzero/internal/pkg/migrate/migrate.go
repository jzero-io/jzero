package migrate

import (
	"fmt"
	"strings"

	"github.com/eddieowens/opts"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	defaultXMigrationsTable = "schema_migrations"
	defaultSource           = "file://desc/sql_migration"
)

type Driver string

const (
	MySQL  Driver = "mysql"
	Pgx    Driver = "pgx"
	Sqlite Driver = "sqlite"
)

type (
	MigrateOpts struct {
		Source             string
		XMigrationsTable   string
		SourceAppendDriver bool
	}

	Migrate struct {
		migrate *migrate.Migrate
	}
)

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
		Source:           defaultSource,
		XMigrationsTable: defaultXMigrationsTable,
	}
}

func NewMigrate(conf sqlx.SqlConf, op ...opts.Opt[MigrateOpts]) (*Migrate, error) {
	ops := opts.DefaultApply(op...)

	var (
		dataSource     = conf.DataSource
		source         = ops.Source
		paramConnector string
	)

	if strings.Contains(dataSource, "?") {
		paramConnector = "&"
	} else {
		paramConnector = "?"
	}

	switch Driver(conf.DriverName) {
	case MySQL:
		dataSource = "mysql://" + dataSource
	case Pgx:
		dataSource = "pgx5://" + strings.TrimPrefix(dataSource, "postgres://")
	case Sqlite:
		dataSource = "sqlite://" + dataSource
	default:
		return nil, fmt.Errorf("unsupported driver: %s", conf.DriverName)
	}

	if ops.SourceAppendDriver {
		source = fmt.Sprintf("%s/%s", source, conf.DriverName)
	}

	dataSource = fmt.Sprintf("%s%sx-migrations-table=%s", dataSource, paramConnector, ops.XMigrationsTable)

	m, err := migrate.New(source, dataSource)
	if err != nil {
		return nil, err
	}

	return &Migrate{
		migrate: m,
	}, nil
}

func (d *Migrate) Up(steps ...uint) error {
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

func (d *Migrate) Down(steps ...uint) error {
	if len(steps) > 1 {
		return errors.New("steps number should not be more than 1")
	}

	if len(steps) == 0 {
		return d.migrate.Down()
	}

	return d.migrate.Steps(-int(steps[0]))
}

func (d *Migrate) Goto(version uint) error {
	return d.migrate.Migrate(version)
}

func (d *Migrate) Version() (uint, bool, error) {
	return d.migrate.Version()
}
