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

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	jzerodesc "github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/gitstatus"
	"github.com/jzero-io/jzero/pkg/osx"
)

type JzeroModel struct {
	Wd     string
	Style  string
	Module string

	config.GenConfig
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
	)

	if jm.ModelMysqlDatasource {
		allTables, err = getMysqlAllTables(jm.ModelMysqlDatasourceUrl)
		if err != nil {
			return err
		}

		writeTables, err := jm.GenDDL(allTables)
		if err != nil {
			return err
		}
		if !jm.GenMysqlCreateTableDDL {
			defer func() {
				for _, v := range writeTables {
					if err = os.Remove(v); err != nil {
						logx.Debugf("remove write ddl file error: %s", err.Error())
					}
				}
			}()
		}
	}

	sqlDir := filepath.Join("desc", "sql")
	if !pathx.FileExists(sqlDir) {
		return nil
	}

	var (
		allFiles        []string
		genCodeSqlFiles []string
	)
	genCodeSqlSpecMap := make(map[string][]*parser.Table)

	allFiles, err = jzerodesc.FindSqlFiles(sqlDir)
	if err != nil {
		return err
	}

	switch {
	case jm.GitChange && len(jm.Desc) == 0:
		if jm.ModelMysqlDatasource {
			// 从 struct migrate 而来
			m, _, err := gitstatus.ChangedFiles(jm.ModelGitChangePath, ".go")
			if err == nil {
				for _, v := range m {
					genCodeTables = append(genCodeTables, getTableNameByGoMethod(v)...)
				}
				for _, v := range genCodeTables {
					genCodeSqlFiles = append(genCodeSqlFiles, filepath.Join("desc", "sql", v+".sql"))
				}
			}
		} else {
			m, _, err := gitstatus.ChangedFiles("desc", ".sql")
			if err == nil {
				genCodeSqlFiles = append(genCodeSqlFiles, m...)
			}
		}
	case len(jm.Desc) > 0:
		for _, v := range jm.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".sql" {
					genCodeSqlFiles = append(genCodeSqlFiles, v)
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
		genCodeSqlFiles, err = jzerodesc.FindSqlFiles(sqlDir)
		if err != nil {
			return err
		}
	}

	// ignore sql desc
	for _, v := range jm.DescIgnore {
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

	if len(genCodeSqlFiles) != 0 {
		var eg errgroup.Group
		for _, f := range allFiles {
			eg.Go(func() error {
				tableParsers, err := parser.Parse(filepath.Join(jm.Wd, f), "", jm.ModelMysqlStrict)
				if err != nil {
					return err
				}
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
		cmd := exec.Command("goctl", "model", "mysql", "ddl", "--database", jm.ModelMysqlDDLDatabase, "--src", f, "--dir", modelDir, "--home", goctlHome, "--style", jm.Style, "-i", strings.Join(jm.ModelMysqlIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", jm.ModelMysqlCache), "--strict="+fmt.Sprintf("%t", jm.ModelMysqlStrict))
		resp, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
		}
		if jm.ModelMysqlCachePrefix != "" && jm.ModelMysqlCache {
			for _, tp := range tableParsers {
				namingFormat, err := format.FileNamingFormat(jm.Style, tp.Name.Source())
				if err != nil {
					return err
				}
				file := namingFormat + "model_gen.go"
				if jm.Style == "go_zero" {
					file = namingFormat + "_model_gen.go"
				}
				err = jm.addModelMysqlCachePrefix(filepath.Join(modelDir, file))
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

func (jm *JzeroModel) addModelMysqlCachePrefix(fp string) error {
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
								basicLit.Value = fmt.Sprintf(`"%s%s"`, jm.ModelMysqlCachePrefix, strings.ReplaceAll(basicLit.Value, "\"", ""))
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

type Table struct {
	Name string `db:"name"`
}

func getMysqlAllTables(url string) ([]string, error) {
	sqlConn := sqlx.NewSqlConn("mysql", url)

	var tables []string
	err := sqlConn.QueryRowsCtx(context.Background(), &tables, "show tables")
	if err != nil {
		return nil, err
	}
	return tables, nil
}

type ShowCreateTableResult struct {
	DDL string `db:"Create Table"`
}

func getTableDDL(url, table string) (string, error) {
	sqlConn := sqlx.NewSqlConn("mysql", url)

	var showCreateTableResult ShowCreateTableResult
	err := sqlConn.QueryRowCtx(context.Background(), &showCreateTableResult, "show create table "+table)
	if err != nil {
		return "", err
	}
	return showCreateTableResult.DDL, nil
}

func getTableNameByGoMethod(fp string) []string {
	var tables []string
	if filepath.Ext(fp) == ".go" {
		fset := token.NewFileSet()

		f, err := goparser.ParseFile(fset, fp, nil, goparser.ParseComments)
		if err != nil {
			return nil
		}

		for _, decl := range f.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == "TableName" {
					for _, stmt := range funcDecl.Body.List {
						if retStmt, ok := stmt.(*ast.ReturnStmt); ok {
							if len(retStmt.Results) > 0 {
								if basicLit, ok := retStmt.Results[0].(*ast.BasicLit); ok {
									tables = append(tables, strings.ReplaceAll(basicLit.Value, `"`, ""))
								}
							}
						}
					}
				}
			}
		}
	}
	return tables
}
