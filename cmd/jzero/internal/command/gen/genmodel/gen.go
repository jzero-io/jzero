package genmodel

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/dsn"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/gitstatus"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
)

type JzeroModel struct {
	Module string
}

func (jm *JzeroModel) Gen() error {
	var (
		allTables     []string
		err           error
		genCodeTables []string
		sqlConn       sqlx.SqlConn
	)

	if config.C.Gen.ModelDriver == "postgres" && !config.C.Gen.ModelDatasource {
		return errors.New("postgres model only support datasource mode")
	}

	if config.C.Gen.ModelDatasource {
		switch config.C.Gen.ModelDriver {
		case "mysql":
			sqlConn = sqlx.NewMysql(config.C.Gen.ModelDatasourceUrl)
		case "postgres":
			sqlConn = postgres.New(config.C.Gen.ModelDatasourceUrl)
		default:
			return errors.Errorf("model driver %s not support", config.C.Gen.ModelDriver)
		}

		tables, err := getAllTables(sqlConn, config.C.Gen.ModelDriver)
		if err != nil {
			return err
		}

		fmt.Printf("%s to generate ddl from %s\n", color.WithColor("Start", color.FgGreen), config.C.Gen.ModelDatasourceUrl)

		writeTables, err := jm.GenDDL(sqlConn, tables)
		if err != nil {
			return err
		}
		if !config.C.Gen.ModelCreateTableDDL {
			defer func() {
				for _, v := range writeTables {
					if err = os.Remove(v); err != nil {
						logx.Debugf("remove write ddl file error: %s", err.Error())
					}
				}
			}()
		}
	}

	if !pathx.FileExists(config.C.SqlDir()) {
		return nil
	}

	var goctlHome string

	if !pathx.FileExists(filepath.Join(config.C.Gen.Home, "go-zero", "model")) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
		err = embeded.WriteTemplateDir(filepath.Join("go-zero", "model"), filepath.Join(tempDir, "model"))
		if err != nil {
			return err
		}
		goctlHome = tempDir
	} else {
		goctlHome = filepath.Join(config.C.Gen.Home, "go-zero")
	}
	logx.Debugf("goctl_home = %s", goctlHome)

	var (
		allFiles        []string
		genCodeSqlFiles []string
	)
	genCodeSqlSpecMap := make(map[string][]*parser.Table)

	allFiles, err = jzerodesc.FindSqlFiles(config.C.SqlDir())
	if err != nil {
		return err
	}

	switch {
	case config.C.Gen.GitChange && filex.DirExists(filepath.Join(config.C.Wd(), ".git")) && len(config.C.Gen.Desc) == 0 && !config.C.Gen.ModelDatasource:
		m, _, err := gitstatus.ChangedFiles(config.C.SqlDir(), ".sql")
		if err == nil {
			genCodeSqlFiles = append(genCodeSqlFiles, m...)
		}
	case len(config.C.Gen.Desc) > 0:
		for _, v := range config.C.Gen.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".sql" {
					genCodeSqlFiles = append(genCodeSqlFiles, filepath.Clean(v))
				}
			} else {
				specifiedSqlFiles, err := jzerodesc.FindSqlFiles(v)
				if err != nil {
					return err
				}
				genCodeSqlFiles = append(genCodeSqlFiles, specifiedSqlFiles...)
			}
		}
	default:
		genCodeSqlFiles, err = jzerodesc.FindSqlFiles(config.C.SqlDir())
		if err != nil {
			return err
		}
	}

	// ignore sql desc
	for _, v := range config.C.Gen.DescIgnore {
		if !osx.IsDir(v) {
			if filepath.Ext(v) == ".sql" {
				genCodeSqlFiles = lo.Reject(genCodeSqlFiles, func(item string, _ int) bool {
					return item == v
				})
			}
		} else {
			specifiedSqlFiles, err := jzerodesc.FindSqlFiles(v)
			if err != nil {
				return err
			}
			for _, saf := range specifiedSqlFiles {
				genCodeSqlFiles = lo.Reject(genCodeSqlFiles, func(item string, _ int) bool {
					return item == saf
				})
			}
		}
	}

	var mu sync.Mutex

	if len(genCodeSqlFiles) != 0 {
		if config.C.Gen.ModelDatasource {
			tables, err := getAllTables(sqlConn, config.C.Gen.ModelDriver)
			if err != nil {
				return err
			}
			if len(config.C.Gen.ModelDatasourceTable) != 0 && config.C.Gen.ModelDatasourceTable[0] != "*" {
				for _, v := range tables {
					if lo.Contains(config.C.Gen.ModelDatasourceTable, cast.ToString(v)) {
						allTables = append(allTables, v)
					}
				}
			} else if len(config.C.Gen.ModelDatasourceTable) != 0 && config.C.Gen.ModelDatasourceTable[0] == "*" {
				allTables = tables
			}
			for _, f := range allFiles {
				genCodeSqlSpecMap[f] = []*parser.Table{
					{
						Name: stringx.From(filepath.Base(f)),
					},
				}
			}
		} else {
			var eg errgroup.Group
			for _, f := range allFiles {
				eg.Go(func() error {
					tableParsers, err := parser.Parse(filepath.Join(config.C.Wd(), f), "", config.C.Gen.ModelStrict)
					if err != nil {
						return err
					}
					mu.Lock()
					defer mu.Unlock()
					genCodeSqlSpecMap[f] = tableParsers
					for _, tp := range tableParsers {
						allTables = append(allTables, tp.Name.Source())
					}
					return nil
				})
			}
			if err = eg.Wait(); err != nil {
				return err
			}
		}
	} else {
		return nil
	}

	fmt.Printf("%s to generate model code from sql files.\n", color.WithColor("Start", color.FgGreen))

	var eg errgroup.Group
	eg.SetLimit(len(genCodeSqlFiles))
	for _, f := range genCodeSqlFiles {
		eg.Go(func() error {
			fmt.Printf("%s sql file %s\n", color.WithColor("Using", color.FgGreen), f)
			tableParsers := genCodeSqlSpecMap[f]

			for _, tp := range tableParsers {
				genCodeTables = append(genCodeTables, tp.Name.Source())
			}

			bf := filepath.Base(f)
			modelDir := filepath.Join("internal", "model", strings.ToLower(bf[0:len(bf)-len(path.Ext(bf))]))

			var scheme string
			if config.C.Gen.ModelDatasource && config.C.Gen.ModelDriver == "mysql" {
				meta, err := dsn.ParseDSN(config.C.Gen.ModelDriver, config.C.Gen.ModelDatasourceUrl)
				if err != nil {
					return err
				}
				scheme = meta[dsn.Database]
			}
			if config.C.Gen.ModelScheme != "" {
				scheme = config.C.Gen.ModelScheme
			}

			if config.C.Gen.ModelDriver == "postgres" {
				if scheme == "" {
					scheme = "public"
				}
				cmd := exec.Command("goctl", "model", "pg", "datasource", "--url", config.C.Gen.ModelDatasourceUrl, "--scheme", scheme, "-t", strings.TrimSuffix(filepath.Base(f), ".sql"), "--dir", modelDir, "--home", goctlHome, "--style", config.C.Gen.Style, "-i", strings.Join(config.C.Gen.ModelIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", config.C.Gen.ModelCache), "-p", config.C.Gen.ModelCachePrefix, "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelStrict))
				logx.Debug(cmd.String())
				resp, err := cmd.CombinedOutput()
				if err != nil {
					return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
				}
			} else {
				cmd := exec.Command("goctl", "model", "mysql", "ddl", "--database", scheme, "--src", f, "--dir", modelDir, "--home", goctlHome, "--style", config.C.Gen.Style, "-i", strings.Join(config.C.Gen.ModelIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", config.C.Gen.ModelCache), "-p", config.C.Gen.ModelCachePrefix, "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelStrict))
				logx.Debug(cmd.String())
				resp, err := cmd.CombinedOutput()
				if err != nil {
					return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
				}
			}
			return nil
		})
	}

	if err = eg.Wait(); err != nil {
		return err
	}

	err = jm.GenRegister(allTables)
	if err != nil {
		return err
	}

	fmt.Println(color.WithColor("Done", color.FgGreen))

	return nil
}

func getAllTables(sqlConn sqlx.SqlConn, driver string) ([]string, error) {
	var tables []string

	switch driver {
	case "mysql":
		err := sqlConn.QueryRowsCtx(context.Background(), &tables, "show tables")
		if err != nil {
			return nil, err
		}
	case "postgres":
		err := sqlConn.QueryRowsCtx(context.Background(), &tables, "select tablename from pg_tables where schemaname = 'public'")
		if err != nil {
			return nil, err
		}
	}
	return tables, nil
}

type ShowCreateTableResult struct {
	DDL string `db:"Create Table"`
}

func getTableDDL(sqlConn sqlx.SqlConn, driver, table string) (string, error) {
	if driver == "postgres" {
		return "-- todo", nil
	}

	var showCreateTableResult ShowCreateTableResult
	err := sqlConn.QueryRowCtx(context.Background(), &showCreateTableResult, "show create table "+table)
	if err != nil {
		return "", err
	}
	return showCreateTableResult.DDL, nil
}
