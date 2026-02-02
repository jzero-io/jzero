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
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/dsn"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
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

		if !config.C.Quiet {
			fmt.Printf("%s to generate model from %s\n", console.Green("Start"), config.C.Gen.ModelDatasourceUrl)
		}
	}

	if !pathx.FileExists(config.C.SqlDir()) && !config.C.Gen.ModelDatasource {
		return nil
	}

	// 处理模板
	var goctlHome string
	tempDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 先写入内置模板
	err = embeded.WriteTemplateDir(filepath.Join("go-zero", "model"), filepath.Join(tempDir, "model"))
	if err != nil {
		return err
	}

	// 如果用户自定义了模板，则复制覆盖
	customTemplatePath := filepath.Join(config.C.Home, "go-zero", "model")
	if pathx.FileExists(customTemplatePath) {
		err = filex.CopyDir(customTemplatePath, filepath.Join(tempDir, "model"))
		if err != nil {
			return err
		}
	}

	goctlHome = tempDir
	logx.Debugf("goctl_home = %s", goctlHome)

	var (
		sqlFiles        []string
		genCodeSqlFiles []string
	)
	genCodeSqlSpecMap := make(map[string][]*ddlparser.Table)

	if !config.C.Gen.ModelDatasource {
		sqlFiles, err = jzerodesc.FindSqlFiles(config.C.SqlDir())
		if err != nil {
			return err
		}

		switch {
		case config.C.Gen.GitChange && gitstatus.IsGitRepo(filepath.Join(config.C.Wd())) && len(config.C.Gen.Desc) == 0:
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
						return item == filepath.Clean(v)
					})
					sqlFiles = lo.Reject(sqlFiles, func(item string, _ int) bool {
						return item == filepath.Clean(v)
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
					sqlFiles = lo.Reject(sqlFiles, func(item string, _ int) bool {
						return item == saf
					})
				}
			}
		}
	}

	var mu sync.Mutex

	if config.C.Gen.ModelDatasource && len(config.C.Gen.Desc) == 0 {
		if len(config.C.Gen.ModelDatasourceTable) == 1 && config.C.Gen.ModelDatasourceTable[0] == "*" {
			allTables, err = getAllTables(conns, config.C.Gen.ModelDriver)
			if err != nil {
				return err
			}
		} else {
			allTables = config.C.Gen.ModelDatasourceTable
		}
		// For datasource mode, generate code for each table directly
		var eg errgroup.Group
		eg.SetLimit(len(allTables))
		for _, tableName := range allTables {
			eg.Go(func() error {
				return generateModelFromDatasource(tableName, goctlHome)
			})
		}
		if err = eg.Wait(); err != nil {
			return err
		}
		return jm.GenRegister(allTables)
	} else if len(genCodeSqlFiles) != 0 {
		var eg errgroup.Group
		for _, f := range sqlFiles {
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
	} else {
		return nil
	}

	if !config.C.Quiet {
		fmt.Printf("%s to generate model code from sql files\n", console.Green("Start"))
	}

	var eg errgroup.Group
	eg.SetLimit(len(genCodeSqlFiles))
	for _, f := range genCodeSqlFiles {
		eg.Go(func() error {
			tableParsers := genCodeSqlSpecMap[f]
			for _, tp := range tableParsers {
				genCodeTables = append(genCodeTables, tp.Name)
			}
			return generateModelFromSqlFile(f, goctlHome)
		})
	}

	if err = eg.Wait(); err != nil {
		return err
	}

	err = jm.GenRegister(allTables)
	if err != nil {
		return err
	}

	if !config.C.Quiet {
		fmt.Println(console.Green("Done"))
	}

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
				allTables = append(allTables, v)
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
				allTables = append(allTables, v)
			}
		}
	}
	return allTables, nil
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

func generateModelFromDatasource(tableName, goctlHome string) error {
	if !config.C.Quiet {
		fmt.Printf("%s table %s from datasource\n", console.Green("Generating"), tableName)
	}

	bf := tableName
	if strings.Contains(tableName, ".") {
		bf = strings.Split(tableName, ".")[1]
	}

	var (
		modelDir string
		schema   = config.C.Gen.ModelSchema
	)

	if strings.Contains(tableName, ".") {
		split := strings.Split(tableName, ".")
		modelDir = filepath.Join("internal", "model", split[0], strings.ToLower(split[1]))
	} else {
		modelDir = filepath.Join("internal", "model", strings.ToLower(bf))
	}

	if config.C.Gen.ModelDriver == "pgx" {
		if schema == "" {
			schema = "public"
		}
		var datasourceUrl string
		if strings.Contains(tableName, ".") {
			for _, v := range config.C.Gen.ModelDatasourceUrl {
				meta, err := dsn.ParseDSN("pgx", v)
				if err != nil {
					return err
				}
				if meta[dsn.Database] == strings.Split(tableName, ".")[0] {
					datasourceUrl = v
					break
				}
			}
		} else {
			datasourceUrl = config.C.Gen.ModelDatasourceUrl[0]
		}

		cmd := exec.Command("goctl", "model", "pg", "datasource", "--url", datasourceUrl, "--schema", schema, "-t", bf, "--dir", modelDir, "--home", goctlHome, "--style", config.C.Style, "-i", strings.Join(getIgnoreColumns(bf), ","), "--cache="+fmt.Sprintf("%t", getIsCacheTable(bf)), "-p", config.C.Gen.ModelCachePrefix, "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelStrict))
		logx.Debug(cmd.String())
		resp, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
		}
	} else if config.C.Gen.ModelDriver == "mysql" {
		var datasourceUrl string
		if strings.Contains(tableName, ".") {
			for _, v := range config.C.Gen.ModelDatasourceUrl {
				meta, err := dsn.ParseDSN("mysql", v)
				if err != nil {
					return err
				}
				if meta[dsn.Database] == strings.Split(tableName, ".")[0] {
					datasourceUrl = v
					break
				}
			}
		} else {
			datasourceUrl = config.C.Gen.ModelDatasourceUrl[0]
		}

		cmd := exec.Command("goctl", "model", "mysql", "datasource", "--url", datasourceUrl, "-t", bf, "--dir", modelDir, "--home", goctlHome, "--style", config.C.Style, "-i", strings.Join(getIgnoreColumns(bf), ","), "--cache="+fmt.Sprintf("%t", getIsCacheTable(bf)), "-p", config.C.Gen.ModelCachePrefix, "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelStrict))
		logx.Debug(cmd.String())
		resp, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
		}
	}

	return nil
}

func generateModelFromSqlFile(sqlFile, goctlHome string) error {
	if !config.C.Quiet {
		fmt.Printf("%s sql file %s\n", console.Green("Using"), sqlFile)
	}

	bf := strings.TrimSuffix(filepath.Base(sqlFile), ".sql")

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

	cmd := exec.Command("goctl", "model", "mysql", "ddl", "--database", schema, "--src", sqlFile, "--dir", modelDir, "--home", goctlHome, "--style", config.C.Style, "-i", strings.Join(getIgnoreColumns(bf), ","), "--cache="+fmt.Sprintf("%t", getIsCacheTable(bf)), "-p", config.C.Gen.ModelCachePrefix, "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelStrict))
	logx.Debug(cmd.String())
	resp, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
	}
	return nil
}

// Parse parses ddl into golang structure
func ParseSql(filename string) ([]*ddlparser.Table, error) {
	p := ddlparser.NewParser()
	tables, err := p.From(filename)
	if err != nil {
		return nil, err
	}

	// Restrict SQL files to only one table
	if len(tables) != 1 {
		return nil, errors.Errorf("SQL file %s contains %d tables, but only one table per SQL file is allowed", filename, len(tables))
	}

	return tables, nil
}
