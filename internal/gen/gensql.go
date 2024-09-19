package gen

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
)

type JzeroSql struct {
	Wd    string
	Style string

	ModelStrict               bool
	ModelIgnoreColumns        []string
	ModelMysqlDatasource      bool
	ModelMysqlDatasourceUrl   string
	ModelMysqlDatasourceTable []string
	ModelMysqlCache           bool
	ModelMysqlCachePrefix     string
}

func (js *JzeroSql) Gen() error {
	dir := "."
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
		// get tables from url
		tables, err := getMysqlAllTables(js.ModelMysqlDatasourceUrl)
		if err != nil {
			return err
		}

		mr.ForEach(func(source chan<- string) {
			for _, table := range tables {
				source <- table
			}
		}, func(table string) {
			fmt.Printf("%s table %s\n", color.WithColor("Using", color.FgGreen), table)

			cmd := exec.Command("goctl", "model", "mysql", "datasource", "--url", js.ModelMysqlDatasourceUrl, "--table", table, "--dir", filepath.Join(dir, "internal", "model", strings.ToLower(table)), "--home", goctlHome, "--style", js.Style, "-i", strings.Join(js.ModelIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", js.ModelMysqlCache), "--strict="+fmt.Sprintf("%t", js.ModelStrict))
			resp, err := cmd.CombinedOutput()
			if err != nil {
				console.Warning("[warning]: %s:%s", err.Error(), resp)
				return
			}

			if js.ModelMysqlCachePrefix != "" && js.ModelMysqlCache {
				namingFormat, _ := format.FileNamingFormat(table, js.Style)
				err = js.addModelMysqlCachePrefix(filepath.Join(dir, "internal", "model", strings.ToLower(table), namingFormat+"model_gen.go"))
				if err != nil {
					console.Warning("[warning]: %s", err.Error())
					return
				}
			}
		}, mr.WithWorkers(len(tables)))
		fmt.Println(color.WithColor("Done", color.FgGreen))
		return nil
	} else {
		sqlDir := filepath.Join(js.Wd, "desc", "sql")
		if f, err := os.Stat(sqlDir); err == nil && f.IsDir() {
			fs, err := os.ReadDir(sqlDir)
			if err != nil {
				return err
			}
			fmt.Printf("%s to generate model code.\n", color.WithColor("Start", color.FgGreen))
			for _, f := range fs {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
					sqlFilePath := filepath.Join(sqlDir, f.Name())
					fmt.Printf("%s sql file %s\n", color.WithColor("Using", color.FgGreen), sqlFilePath)

					tables, err := parser.Parse(sqlFilePath, "", false)
					if err != nil {
						return err
					}

					modelDir := filepath.Join(dir, "internal", "model", strings.ToLower(f.Name()[0:len(f.Name())-len(path.Ext(f.Name()))]))
					cmd := exec.Command("goctl", "model", "mysql", "ddl", "--src", filepath.Join(dir, "desc", "sql", f.Name()), "--dir", modelDir, "--home", goctlHome, "--style", js.Style, "-i", strings.Join(js.ModelIgnoreColumns, ","), "--cache="+fmt.Sprintf("%t", js.ModelMysqlCache), "--strict="+fmt.Sprintf("%t", js.ModelStrict))
					resp, err := cmd.CombinedOutput()
					if err != nil {
						console.Warning("[warning]: %s:%s", err.Error(), resp)
						continue
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
			}
			fmt.Println(color.WithColor("Done", color.FgGreen))
		}
	}
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
	// Write the modified AST back to the file
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
