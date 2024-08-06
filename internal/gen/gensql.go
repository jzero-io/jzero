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
	"path"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/embeded"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

type JzeroSql struct {
	Wd                        string
	Style                     string
	ModelIgnoreColumns        []string
	ModelMysqlDatasource      bool
	ModelMysqlDatasourceUrl   string
	ModelMysqlDatasourceTable []string
	ModelMysqlCache           bool
	ModelMysqlCachePrefix     string
}

func (js *JzeroSql) Gen() error {
	dir := "."

	if js.ModelMysqlDatasource {
		fmt.Printf("%s to generate model code from url %s.\n", color.WithColor("Start", color.FgGreen), js.ModelMysqlDatasourceUrl)
		// get tables from url
		tables, err := getMysqlAllTables(js.ModelMysqlDatasourceUrl)
		if err != nil {
			return err
		}

		for _, table := range tables {
			fmt.Printf("%s table %s\n", color.WithColor("Using", color.FgGreen), table)
			command := fmt.Sprintf("goctl model mysql datasource --url '%s' --table %s --dir %s --home %s --style %s -i '%s' --cache=%t",
				js.ModelMysqlDatasourceUrl,
				table,
				filepath.Join(dir, "internal", "model", strings.ToLower(table)),
				filepath.Join(embeded.Home, "go-zero"),
				js.Style,
				strings.Join(js.ModelIgnoreColumns, ","),
				js.ModelMysqlCache,
			)
			_, err := execx.Run(command, js.Wd)
			if err != nil {
				console.Warning("[warning]: %s", err.Error())
				continue
			}

			if js.ModelMysqlCachePrefix != "" && js.ModelMysqlCache {
				namingFormat, err := format.FileNamingFormat(table, js.Style)
				if err != nil {
					return err
				}
				err = js.addModelMysqlCachePrefix(filepath.Join(dir, "internal", "model", strings.ToLower(table), namingFormat+"model_gen.go"))
				if err != nil {
					return err
				}
			}
			fmt.Println(color.WithColor("Done", color.FgGreen))
		}
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
					command := fmt.Sprintf("goctl model mysql ddl --src %s --dir %s --home %s --style %s -i '%s' --cache=%t",
						filepath.Join(dir, "desc", "sql", f.Name()),
						modelDir,
						filepath.Join(embeded.Home, "go-zero"),
						js.Style,
						strings.Join(js.ModelIgnoreColumns, ","),
						js.ModelMysqlCache,
					)
					_, err = execx.Run(command, js.Wd)
					if err != nil {
						console.Warning("[warning]: %s", err.Error())
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
