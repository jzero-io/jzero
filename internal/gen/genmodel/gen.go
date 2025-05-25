package genmodel

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	goformat "go/format"
	goparser "go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/jzero-io/jzero-contrib/filex"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	jzerodesc "github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/gitstatus"
	"github.com/jzero-io/jzero/pkg/osx"
)

type JzeroModel struct {
	Module string
	IsNew  bool
}

func (jm *JzeroModel) Gen() error {
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
		allTables     []string
		err           error
		genCodeTables []string
		sqlConn       sqlx.SqlConn
	)

	if config.C.Gen.ModelDriver == "postgres" && !config.C.Gen.ModelDatasource {
		return errors.New("postgres model only support datasource mode")
	}

	if config.C.Gen.ModelMysqlDatasource || config.C.Gen.ModelDatasource {
		if jm.IsNew {
			fmt.Printf("%s you are using mysql datesource to generate model code, please manual execute jzero gen command\n", color.WithColor("Detected", color.FgRed))
			return nil
		}

		switch config.C.Gen.ModelDriver {
		case "mysql":
			if config.C.Gen.ModelMysqlDatasourceUrl != "" {
				config.C.Gen.ModelDatasourceUrl = config.C.Gen.ModelMysqlDatasourceUrl
			}
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

		fmt.Printf("%s to generate ddl from %s\n", color.WithColor("Start", color.FgGreen), config.C.Gen.ModelMysqlDatasourceUrl)

		writeTables, err := jm.GenDDL(sqlConn, tables)
		if err != nil {
			return err
		}
		if !config.C.Gen.ModelMysqlCreateTableDDL || !config.C.Gen.ModelCreateTableDDL {
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
	case config.C.Gen.GitChange && filex.DirExists(filepath.Join(config.C.Wd(), ".git")) && len(config.C.Gen.Desc) == 0 && !config.C.Gen.ModelMysqlDatasource:
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
		if config.C.Gen.ModelMysqlDatasource || config.C.Gen.ModelDatasource {
			tables, err := getAllTables(sqlConn, config.C.Gen.ModelDriver)
			if err != nil {
				return err
			}
			if (len(config.C.Gen.ModelMysqlDatasourceTable) != 0 && config.C.Gen.ModelMysqlDatasourceTable[0] != "*") || (len(config.C.Gen.ModelDatasourceTable) != 0 && config.C.Gen.ModelDatasourceTable[0] != "*") {
				for _, v := range tables {
					if lo.Contains(config.C.Gen.ModelMysqlDatasourceTable, cast.ToString(v)) || lo.Contains(config.C.Gen.ModelDatasourceTable, cast.ToString(v)) {
						allTables = append(allTables, v)
					}
				}
			} else if (len(config.C.Gen.ModelMysqlDatasourceTable) != 0 && config.C.Gen.ModelMysqlDatasourceTable[0] == "*") || (len(config.C.Gen.ModelDatasourceTable) != 0 && config.C.Gen.ModelDatasourceTable[0] == "*") {
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
					tableParsers, err := parser.Parse(filepath.Join(config.C.Wd(), f), "", config.C.Gen.ModelMysqlStrict)
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
	for _, f := range genCodeSqlFiles {
		fmt.Printf("%s sql file %s\n", color.WithColor("Using", color.FgGreen), f)
		tableParsers := genCodeSqlSpecMap[f]

		for _, tp := range tableParsers {
			genCodeTables = append(genCodeTables, tp.Name.Source())
		}

		bf := filepath.Base(f)
		modelDir := filepath.Join("internal", "model", strings.ToLower(bf[0:len(bf)-len(path.Ext(bf))]))

		var ddlDatabase string
		if config.C.Gen.ModelMysqlDDLDatabase != "" {
			ddlDatabase = config.C.Gen.ModelMysqlDDLDatabase
		} else if config.C.Gen.ModelMysqlDatasourceUrl != "" {
			mysqlDsn, err := mysql.ParseDSN(config.C.Gen.ModelMysqlDatasourceUrl)
			if err != nil {
				return err
			}
			ddlDatabase = mysqlDsn.DBName
		}

		if config.C.Gen.ModelDriver == "postgres" {
			cmd := exec.Command("goctl", "model", "pg", "datasource", "--url", config.C.Gen.ModelDatasourceUrl, "-t", strings.TrimSuffix(filepath.Base(f), ".sql"), "--dir", modelDir, "--home", goctlHome, "--style", config.C.Gen.Style, "-i", strings.Join(config.C.Gen.ModelIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", config.C.Gen.ModelCache), "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelStrict))
			logx.Debug(cmd.String())
			resp, err := cmd.CombinedOutput()
			if err != nil {
				return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
			}
		} else {
			if config.C.Gen.ModelMysqlCache {
				config.C.Gen.ModelCache = true
			}
			cmd := exec.Command("goctl", "model", "mysql", "ddl", "--database", ddlDatabase, "--src", f, "--dir", modelDir, "--home", goctlHome, "--style", config.C.Gen.Style, "-i", strings.Join(config.C.Gen.ModelMysqlIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", config.C.Gen.ModelCache), "--strict="+fmt.Sprintf("%t", config.C.Gen.ModelMysqlStrict))
			logx.Debug(cmd.String())
			resp, err := cmd.CombinedOutput()
			if err != nil {
				return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
			}
		}

		if (config.C.Gen.ModelMysqlCachePrefix != "" && config.C.Gen.ModelMysqlCache) || (config.C.Gen.ModelCachePrefix != "" && config.C.Gen.ModelCache) {
			for _, tp := range tableParsers {
				namingFormat, err := format.FileNamingFormat(config.C.Gen.Style, tp.Name.Source())
				if err != nil {
					return err
				}
				file := namingFormat + "model_gen.go"
				if config.C.Gen.Style == "go_zero" {
					file = namingFormat + "_model_gen.go"
				}
				err = jm.addModelCachePrefix(filepath.Join(modelDir, file))
				if err != nil {
					return err
				}
			}
		}
	}

	err = jm.GenRegister(allTables)
	if err != nil {
		return err
	}

	fmt.Println(color.WithColor("Done", color.FgGreen))

	return nil
}

func (jm *JzeroModel) addModelCachePrefix(fp string) error {
	fset := token.NewFileSet()
	f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(node ast.Node) bool {
		if genDecl, ok := node.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range valueSpec.Names {
						if strings.HasPrefix(name.Name, "cache") && strings.HasSuffix(name.Name, "Prefix") {
							value := valueSpec.Values[i]
							if basicLit, ok := value.(*ast.BasicLit); ok {
								if config.C.Gen.ModelCachePrefix != "" {
									basicLit.Value = fmt.Sprintf(`"%s%s"`, config.C.Gen.ModelCachePrefix, strings.ReplaceAll(basicLit.Value, "\"", ""))
								} else if config.C.Gen.ModelMysqlCachePrefix != "" {
									basicLit.Value = fmt.Sprintf(`"%s%s"`, config.C.Gen.ModelMysqlCachePrefix, strings.ReplaceAll(basicLit.Value, "\"", ""))
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	buf := bytes.NewBuffer(nil)
	if err := goformat.Node(buf, fset, f); err != nil {
		return err
	}

	if err := os.WriteFile(fp, buf.Bytes(), 0o644); err != nil {
		return err
	}
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
