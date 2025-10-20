package genmodel

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	ddlparser "github.com/zeromicro/ddl-parser/parser"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/dsn"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/gitstatus"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
)

type JzeroModel struct {
	Module string
}

type Conn struct {
	Schema  string
	SqlConn sqlx.SqlConn
}

func (jm *JzeroModel) Gen() error {
	var (
		allTables     []string
		err           error
		genCodeTables []string
		conns         []Conn
	)

	if config.C.Gen.ModelDriver == "postgres" {
		config.C.Gen.ModelDriver = "pgx"
	}

	if config.C.Gen.ModelDriver == "pgx" && !config.C.Gen.ModelDatasource {
		return errors.New("postgres model only support datasource mode")
	}

	if config.C.Gen.ModelDatasource {
		switch config.C.Gen.ModelDriver {
		case "mysql":
			for _, v := range config.C.Gen.ModelDatasourceUrl {
				meta, err := dsn.ParseDSN(config.C.Gen.ModelDriver, v)
				if err != nil {
					return err
				}
				conns = append(conns, Conn{
					Schema:  meta[dsn.Database],
					SqlConn: sqlx.NewMysql(v),
				})
			}
		case "pgx":
			for _, v := range config.C.Gen.ModelDatasourceUrl {
				meta, err := dsn.ParseDSN(config.C.Gen.ModelDriver, v)
				if err != nil {
					return err
				}
				conns = append(conns, Conn{
					Schema:  meta[dsn.Database],
					SqlConn: postgres.New(v),
				})
			}
		default:
			return errors.Errorf("model driver %s not support", config.C.Gen.ModelDriver)
		}

		tables, err := getAllTables(conns, config.C.Gen.ModelDriver)
		if err != nil {
			return err
		}

		fmt.Printf("%s to generate ddl from %s\n", color.WithColor("Start", color.FgGreen), config.C.Gen.ModelDatasourceUrl)

		writeTables, err := jm.GenDDL(conns, tables)
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
	genCodeSqlSpecMap := make(map[string][]*ddlparser.Table)

	allFiles, err = jzerodesc.FindSqlFiles(config.C.SqlDir())
	if err != nil {
		return err
	}

	switch {
	case config.C.Gen.GitChange && gitstatus.IsGitRepo(filepath.Join(config.C.Wd())) && len(config.C.Gen.Desc) == 0 && !config.C.Gen.ModelDatasource:
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
			allTables = config.C.Gen.ModelDatasourceTable
			for _, f := range allFiles {
				genCodeSqlSpecMap[f] = []*ddlparser.Table{
					{
						Name: stringx.From(filepath.Base(f)).Source(),
					},
				}
			}
		} else {
			var eg errgroup.Group
			for _, f := range allFiles {
				eg.Go(func() error {
					tableParsers, err := ParseSql(filepath.Join(config.C.Wd(), f))
					if err != nil {
						return err
					}
					mu.Lock()
					defer mu.Unlock()
					genCodeSqlSpecMap[f] = tableParsers

					bf := strings.TrimSuffix(filepath.Base(f), ".sql")

					for _, tp := range tableParsers {
						if strings.Contains(bf, ".") {
							allTables = append(allTables, strings.Split(bf, ".")[0]+"."+tp.Name)
						} else {
							allTables = append(allTables, tp.Name)
						}
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

	fmt.Printf("%s to generate model code from sql files\n", color.WithColor("Start", color.FgGreen))

	var eg errgroup.Group
	eg.SetLimit(len(genCodeSqlFiles))
	for _, f := range genCodeSqlFiles {
		eg.Go(func() error {
			fmt.Printf("%s sql file %s\n", color.WithColor("Using", color.FgGreen), f)
			tableParsers := genCodeSqlSpecMap[f]

			for _, tp := range tableParsers {
				genCodeTables = append(genCodeTables, tp.Name)
			}

			bf := strings.TrimSuffix(filepath.Base(f), ".sql")

			var (
				modelDir string
				schema   = config.C.Gen.ModelSchema
			)
			if strings.Contains(bf, ".") {
				split := strings.Split(bf, ".")
				modelDir = filepath.Join("internal", "model", split[0], strings.ToLower(split[1]))
			} else {
				modelDir = filepath.Join("internal", "model", strings.ToLower(bf))
			}

			if config.C.Gen.ModelDriver == "pgx" {
				if schema == "" {
					schema = "public"
				}
			} else if config.C.Gen.ModelDriver == "mysql" {
				if strings.Contains(bf, ".") {
					schema = strings.Split(bf, ".")[0]
				} else {
					if schema == "" {
						if len(config.C.Gen.ModelDatasourceUrl) >= 1 {
							meta, err := dsn.ParseDSN("mysql", config.C.Gen.ModelDatasourceUrl[0])
							if err != nil {
								return err
							}
							schema = meta[dsn.Database]
						}
					}
				}
			}

			if config.C.Gen.ModelDriver == "pgx" {
				var datasourceUrl string
				if strings.Contains(bf, ".") {
					for _, v := range config.C.Gen.ModelDatasourceUrl {
						meta, err := dsn.ParseDSN("pgx", v)
						if err != nil {
							return err
						}
						if meta[dsn.Database] == strings.Split(bf, ".")[0] {
							datasourceUrl = v
							break
						}
					}
				} else {
					datasourceUrl = config.C.Gen.ModelDatasourceUrl[0]
				}

				tableName := func() string {
					if strings.Contains(bf, ".") {
						return strings.Split(bf, ".")[1]
					}
					return bf
				}()
				cmd := exec.Command("goctl", "model", "pg", "datasource", "--url", datasourceUrl, "--schema", schema, "-t", tableName, "--dir", modelDir, "--home", goctlHome, "--style", config.C.Gen.Style, "-i", strings.Join(getIgnoreColumns(tableName), ","), "--cache="+fmt.Sprintf("%t", getIsCacheTable(bf)), "-p", config.C.Gen.ModelCachePrefix, "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelStrict))
				logx.Debug(cmd.String())
				resp, err := cmd.CombinedOutput()
				if err != nil {
					return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
				}
			} else {
				tableName := bf
				if strings.Contains(bf, ".") {
					tableName = strings.Split(bf, ".")[1]
				}
				cmd := exec.Command("goctl", "model", "mysql", "ddl", "--database", schema, "--src", f, "--dir", modelDir, "--home", goctlHome, "--style", config.C.Gen.Style, "-i", strings.Join(getIgnoreColumns(tableName), ","), "--cache="+fmt.Sprintf("%t", getIsCacheTable(bf)), "-p", config.C.Gen.ModelCachePrefix, "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelStrict))
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

func getAllTables(conns []Conn, driver string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var allTables []string

	switch driver {
	case "mysql":
		for _, conn := range conns {
			var tables []string
			err := conn.SqlConn.QueryRowsCtx(ctx, &tables, "show tables")
			if err != nil {
				return nil, err
			}
			for _, v := range tables {
				allTables = append(allTables, fmt.Sprintf("`%s`", conn.Schema)+"."+fmt.Sprintf("`%s`", v))
			}
		}
	case "pgx":
		if config.C.Gen.ModelSchema == "" {
			config.C.Gen.ModelSchema = "public"
		}
		for _, conn := range conns {
			var tables []string
			err := conn.SqlConn.QueryRowsCtx(ctx, &tables, fmt.Sprintf("select tablename from pg_tables where schemaname = '%s'", config.C.Gen.ModelSchema))
			if err != nil {
				return nil, err
			}
			for _, v := range tables {
				allTables = append(allTables, conn.Schema+"."+v)
			}
		}
	}
	return allTables, nil
}

type ShowCreateTableResult struct {
	DDL string `db:"Create Table"`
}

func getTableDDL(sqlConn sqlx.SqlConn, driver, table string) (string, error) {
	if driver == "pgx" {
		return "-- todo", nil
	}

	var showCreateTableResult ShowCreateTableResult

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := sqlConn.QueryRowCtx(ctx, &showCreateTableResult, "show create table "+table)
	if err != nil {
		return "", err
	}
	return showCreateTableResult.DDL, nil
}

func getIsCacheTable(t string) bool {
	if config.C.Gen.ModelCache && len(config.C.Gen.ModelCacheTable) == 1 && config.C.Gen.ModelCacheTable[0] == "*" {
		return true
	}

	if config.C.Gen.ModelCache {
		for _, v := range config.C.Gen.ModelCacheTable {
			if v == t {
				return true
			}
		}
	}
	return false
}

func getIgnoreColumns(tableName string) []string {
	if config.C.Gen.ModelIgnoreColumnsTable != nil {
		for _, v := range config.C.Gen.ModelIgnoreColumnsTable {
			if v.Table == tableName {
				return v.Columns
			}
		}
	}
	return config.C.Gen.ModelIgnoreColumns
}

// Parse parses ddl into golang structure
func ParseSql(filename string) ([]*ddlparser.Table, error) {
	p := ddlparser.NewParser()
	tables, err := p.From(filename)
	if err != nil {
		return nil, err
	}
	return tables, nil
}
