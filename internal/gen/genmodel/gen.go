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

type JzeroModel struct {
	Wd     string
	Style  string
	Module string

	config.GenConfig
}

func (js *JzeroModel) Gen() error {
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

	if js.ModelMysqlDatasource {
		fmt.Printf("%s to generate model code from url %s.\n", color.WithColor("Start", color.FgGreen), js.ModelMysqlDatasourceUrl)

		allTables, err = getMysqlAllTables(js.ModelMysqlDatasourceUrl)
		if err != nil {
			return err
		}

		if js.GitChange {
			var files []string
			m, _, err := gitstatus.ChangedFiles(js.ModelGitChangePath, "")
			if err == nil {
				files = append(files, m...)
			}

			for _, v := range files {
				genCodeTables = append(genCodeTables, getTableNameByGoMethod(v)...)
			}
		} else {
			genCodeTables = append(genCodeTables, allTables...)
		}

		mr.ForEach(func(source chan<- string) {
			for _, table := range genCodeTables {
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
		}, mr.WithWorkers(len(genCodeTables)))
	} else {
		sqlDir := filepath.Join("desc", "sql")
		if !pathx.FileExists(sqlDir) {
			return nil
		}

		var (
			allFiles     []string
			genCodeFiles []string
		)

		allFiles, err = jzerodesc.FindSqlFiles(sqlDir)
		if err != nil {
			return err
		}

		for _, f := range allFiles {
			tableParsers, err := parser.Parse(filepath.Join(js.Wd, f), "", js.ModelMysqlStrict)
			if err != nil {
				return err
			}
			for _, tp := range tableParsers {
				allTables = append(allTables, tp.Name.Source())
			}
		}

		switch {
		case js.GitChange && len(js.Desc) == 0:
			m, _, err := gitstatus.ChangedFiles(js.ModelGitChangePath, ".sql")
			if err == nil {
				genCodeFiles = append(genCodeFiles, m...)
			}
		case len(js.Desc) > 0:
			for _, v := range js.Desc {
				if !osx.IsDir(v) {
					if filepath.Ext(v) == ".sql" {
						genCodeFiles = append(genCodeFiles, v)
					}
				} else {
					specifiedSqlFiles, err := jzerodesc.FindApiFiles(v)
					if err != nil {
						return err
					}
					genCodeFiles = append(genCodeFiles, specifiedSqlFiles...)
				}
			}
		default:
			genCodeFiles, err = jzerodesc.FindSqlFiles(sqlDir)
			if err != nil {
				return err
			}
		}

		fmt.Printf("%s to generate model code from sql files.\n", color.WithColor("Start", color.FgGreen))
		for _, f := range genCodeFiles {
			fmt.Printf("%s sql file %s\n", color.WithColor("Using", color.FgGreen), f)
			tableParsers, err := parser.Parse(filepath.Join(js.Wd, f), "", js.ModelMysqlStrict)
			if err != nil {
				return err
			}

			for _, tp := range tableParsers {
				genCodeTables = append(genCodeTables, tp.Name.Source())
			}

			bf := filepath.Base(f)
			modelDir := filepath.Join("internal", "model", strings.ToLower(bf[0:len(bf)-len(path.Ext(bf))]))
			cmd := exec.Command("goctl", "model", "mysql", "ddl", "--src", f, "--dir", modelDir, "--home", goctlHome, "--style", js.Style, "-i", strings.Join(js.ModelMysqlIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", js.ModelMysqlCache), "--strict="+fmt.Sprintf("%t", js.ModelMysqlStrict))
			resp, err := cmd.CombinedOutput()
			if err != nil {
				return errors.Errorf("gen model code meet error. Err: %s:%s", err.Error(), resp)
			}
			if js.ModelMysqlCachePrefix != "" && js.ModelMysqlCache {
				for _, tp := range tableParsers {
					namingFormat, err := format.FileNamingFormat(js.Style, tp.Name.Source())
					if err != nil {
						return err
					}
					file := namingFormat + "model_gen.go"
					if js.Style == "go_zero" {
						file = namingFormat + "_model_gen.go"
					}
					err = js.addModelMysqlCachePrefix(filepath.Join(modelDir, file))
					if err != nil {
						return err
					}
				}
			}
		}
	}

	err = js.GenRegister(allTables)
	if err != nil {
		return err
	}
	fmt.Println(color.WithColor("Done", color.FgGreen))

	return nil
}

func (js *JzeroModel) addModelMysqlCachePrefix(fp string) error {
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
