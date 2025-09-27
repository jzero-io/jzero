package migrate

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/eddieowens/opts"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type MigrateOpts struct {
	PreProcessSqlFunc func(content string) string
	Source            string
}

func (opts MigrateOpts) DefaultOptions() MigrateOpts {
	return MigrateOpts{
		PreProcessSqlFunc: func(content string) string {
			return content
		},
		Source: "file://desc/sql_migration",
	}
}

func WithPreProcessSqlFunc(f func(string) string) opts.Opt[MigrateOpts] {
	return func(opts *MigrateOpts) {
		opts.PreProcessSqlFunc = f
	}
}

func WithSource(source string) opts.Opt[MigrateOpts] {
	return func(opts *MigrateOpts) {
		opts.Source = source
	}
}

func Migrate(ctx context.Context, c sqlx.SqlConf, op ...opts.Opt[MigrateOpts]) error {
	ops := opts.DefaultApply(op...)
	var databaseUrl string
	switch c.DriverName {
	case "mysql":
		databaseUrl = "mysql://" + c.DataSource
	case "pgx":
		databaseUrl = "pgx5://" + strings.TrimPrefix(c.DataSource, "postgres://")
	}
	if err := sqlMigrate(ops.Source, databaseUrl, c, ops); err != nil {
		return err
	}
	return nil
}

type customFileSource struct {
	*file.File
	driverName        string
	preProcessSqlFunc func(content string) string
}

func (c *customFileSource) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	rc, id, err := c.File.ReadUp(version)
	if err != nil {
		return nil, "", err
	}

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, "", err
	}

	if err = rc.Close(); err != nil {
		return nil, "", err
	}
	return io.NopCloser(strings.NewReader(c.preProcessSqlFunc(string(content)))), id, nil
}

func (c *customFileSource) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	rc, id, err := c.File.ReadDown(version)
	if err != nil {
		return nil, "", err
	}

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, "", err
	}

	if err = rc.Close(); err != nil {
		return nil, "", err
	}

	modifiedContent := c.preProcessSqlFunc(string(content))
	return io.NopCloser(strings.NewReader(modifiedContent)), id, nil
}

func sqlMigrate(sourceUrl, databaseUrl string, c sqlx.SqlConf, ops MigrateOpts) error {
	fileDriver := &file.File{}
	fileSource, err := fileDriver.Open(sourceUrl)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	customSource := &customFileSource{
		File:              fileSource.(*file.File),
		driverName:        c.DriverName,
		preProcessSqlFunc: ops.PreProcessSqlFunc,
	}

	m, err := migrate.NewWithSourceInstance("file", customSource, databaseUrl)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
	}

	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		return sourceErr
	}
	if databaseErr != nil {
		return databaseErr
	}
	return nil
}
