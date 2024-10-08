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
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	jzerodesc "github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/gitstatus"
	"github.com/jzero-io/jzero/pkg/osx"
)

type JzeroSql struct {
	Wd    string
	Style string
	config.GenConfig
}

func (js *JzeroSql) Gen() error {
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

	if js.ModelMysqlDatasource {
		fmt.Printf("%s to generate model code from url %s.\n", color.WithColor("Start", color.FgGreen), js.ModelMysqlDatasourceUrl)
		tables, err := getMysqlAllTables(js.ModelMysqlDatasourceUrl)
		if err != nil {
			return err
		}
		if js.GitDiff {
			var changesTables []string
			var files []string
			m, _, err := gitstatus.ChangedFiles(js.ModelGitDiffPath, "")
			if err != nil {
				return err
			}
			files = append(files, m...)

			for _, v := range files {
				changesTables = append(changesTables, getTableNameByGoMethod(v)...)
			}
			tables = changesTables
		}

		mr.ForEach(func(source chan<- string) {
			for _, table := range tables {
				source <- table
			}
		}, func(table string) {
			fmt.Printf("%s table %s\n", color.WithColor("Using", color.FgGreen), table)

			cmd := exec.Command("goctl", "model", "mysql", "datasource", "--url", js.ModelMysqlDatasourceUrl, "--table", table, "--dir", filepath.Join("internal", "model", strings.ToLower(table)), "--home", goctlHome, "--style", js.Style, "-i", strings.Join(js.ModelMysqlIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", js.ModelMysqlCache), "--strict="+fmt.Sprintf("%t", js.ModelMysqlStrict))
			resp, err := cmd.CombinedOutput()
			if err != nil {
				console.Warning("[warning]: %s:%s", err.Error(), resp)
				return
			}

			if js.ModelMysqlCachePrefix != "" && js.ModelMysqlCache {
				namingFormat, _ := format.FileNamingFormat(table, js.Style)
				err = js.addModelMysqlCachePrefix(filepath.Join("internal", "model", strings.ToLower(table), namingFormat+"model_gen.go"))
				if err != nil {
					console.Warning("[warning]: %s", err.Error())
					return
				}
			}
		}, mr.WithWorkers(len(tables)))
		fmt.Println(color.WithColor("Done", color.FgGreen))
		return nil
	}

	sqlDir := filepath.Join(js.Wd, "desc", "sql")
	if !pathx.FileExists(sqlDir) {
		return nil
	}

	var (
		files []string
		err   error
	)

	switch {
	case js.GitDiff && len(js.Desc) == 0:
		m, _, err := gitstatus.ChangedFiles(js.ModelGitDiffPath, ".sql")
		if err == nil {
			files = append(files, m...)
		}
	case len(js.Desc) > 0:
		for _, v := range js.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".sql" {
					files = append(files, v)
				}
			} else {
				specifiedSqlFiles, err := jzerodesc.FindApiFiles(v)
				if err != nil {
					return err
				}
				files = append(files, specifiedSqlFiles...)
			}
		}
	default:
		files, err = jzerodesc.FindSqlFiles(sqlDir)
		if err != nil {
			return err
		}
	}

	if len(files) == 0 {
		return nil
	}

	fmt.Printf("%s to generate model code from sql files.\n", color.WithColor("Start", color.FgGreen))
	for _, f := range files {
		fmt.Printf("%s sql file %s\n", color.WithColor("Using", color.FgGreen), f)
		tables, err := parser.Parse(filepath.Join(js.Wd, f), "", js.ModelMysqlStrict)
		if err != nil {
			return err
		}
		bf := filepath.Base(f)
		modelDir := filepath.Join("internal", "model", strings.ToLower(bf[0:len(bf)-len(path.Ext(bf))]))
		cmd := exec.Command("goctl", "model", "mysql", "ddl", "--src", f, "--dir", modelDir, "--home", goctlHome, "--style", js.Style, "-i", strings.Join(js.ModelMysqlIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", js.ModelMysqlCache), "--strict="+fmt.Sprintf("%t", js.ModelMysqlStrict))
		resp, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
		}
		if js.ModelMysqlCachePrefix != "" && js.ModelMysqlCache {
			for _, table := range tables {
				namingFormat, err := format.FileNamingFormat(js.Style, table.Name.Source())
				if err != nil {
					return err
				}
				err = js.addModelMysqlCachePrefix(filepath.Join(modelDir, namingFormat+"model_gen.go"))
				if err != nil {
					return err
				}
			}
		}
	}
	fmt.Println(color.WithColor("Done", color.FgGreen))

	return nil
}

func (js *JzeroSql) addModelMysqlCachePrefix(fp string) error {
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
								basicLit.Value = fmt.Sprintf(`"%s%s"`, js.ModelMysqlCachePrefix, strings.ReplaceAll(basicLit.Value, "\"", ""))
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
