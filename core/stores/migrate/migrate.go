package migrate

import (
	"context"
	"io"
	"strings"

	"github.com/eddieowens/opts"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type MigrateUpOpts struct {
	PreProcessSqlFunc    func(version uint, content string) string
	BeforeMigrateUpFuncs map[uint]func(version uint) error
	AfterMigrateUpFuncs  map[uint]func(version uint) error
	Source               string
}

func (opts MigrateUpOpts) DefaultOptions() MigrateUpOpts {
	return MigrateUpOpts{
		Source: "file://desc/sql_migration",
	}
}

func WithPreProcessSqlFunc(f func(uint, string) string) opts.Opt[MigrateUpOpts] {
	return func(opts *MigrateUpOpts) {
		opts.PreProcessSqlFunc = f
	}
}

func WithSource(source string) opts.Opt[MigrateUpOpts] {
	return func(opts *MigrateUpOpts) {
		opts.Source = source
	}
}

func WithBeforeMigrateUpFunc(mapFuncs map[uint]func(uint) error) opts.Opt[MigrateUpOpts] {
	return func(opts *MigrateUpOpts) {
		opts.BeforeMigrateUpFuncs = mapFuncs
	}
}

func WithAfterMigrateUpFunc(mapFunc map[uint]func(uint) error) opts.Opt[MigrateUpOpts] {
	return func(opts *MigrateUpOpts) {
		opts.AfterMigrateUpFuncs = mapFunc
	}
}

func MigrateUp(ctx context.Context, c sqlx.SqlConf, op ...opts.Opt[MigrateUpOpts]) error {
	ops := opts.DefaultApply(op...)
	var databaseUrl string
	switch c.DriverName {
	case "mysql":
		databaseUrl = "mysql://" + c.DataSource
	case "pgx":
		databaseUrl = "pgx5://" + strings.TrimPrefix(c.DataSource, "postgres://")
	}
	if err := migrateUp(ctx, ops.Source, databaseUrl, c, ops); err != nil {
		return err
	}
	return nil
}

type customFileSource struct {
	*file.File
	preProcessSqlFunc func(version uint, content string) string
	sqlConf           sqlx.SqlConf
	ctx               context.Context
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

	return io.NopCloser(strings.NewReader(c.preProcessSqlFunc(version, string(content)))), id, nil
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

	return io.NopCloser(strings.NewReader(c.preProcessSqlFunc(version, string(content)))), id, nil
}

func migrateUp(ctx context.Context, sourceUrl, databaseUrl string, c sqlx.SqlConf, ops MigrateUpOpts) error {
	fileDriver := &file.File{}
	fileSource, err := fileDriver.Open(sourceUrl)
	if err != nil {
		return err
	}

	customSource := &customFileSource{
		File:              fileSource.(*file.File),
		preProcessSqlFunc: ops.PreProcessSqlFunc,
		sqlConf:           c,
		ctx:               ctx,
	}

	var m *migrate.Migrate
	if ops.PreProcessSqlFunc == nil {
		m, err = migrate.New(sourceUrl, databaseUrl)
		if err != nil {
			return err
		}
	} else {
		m, err = migrate.NewWithSourceInstance("file", customSource, databaseUrl)
		if err != nil {
			return err
		}
	}
	defer m.Close()

	if ops.BeforeMigrateUpFuncs == nil && ops.AfterMigrateUpFuncs == nil {
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		return nil
	}
	// 获取当前版本
	currentVersion, _, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			// 不存在的话, 直接返回 Up
			return m.Up()
		}
	}

	for {
		nextVersion, err := customSource.Next(currentVersion)
		if err == nil && nextVersion > currentVersion {
			if f, ok := ops.BeforeMigrateUpFuncs[nextVersion]; ok {
				if err = f(nextVersion); err != nil {
					return err
				}
			}
			if err = m.Steps(1); err != nil {
				return err
			}
			if f, ok := ops.AfterMigrateUpFuncs[nextVersion]; ok {
				if err = f(nextVersion); err != nil {
					if stepDownErr := m.Steps(-1); stepDownErr != nil {
						return stepDownErr
					}
					return err
				}
			}
			currentVersion = nextVersion
		} else {
			break
		}
	}

	return nil
}
