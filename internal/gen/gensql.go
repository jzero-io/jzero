package gen

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/jzero-io/jzero/embeded"

	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

type JzeroSql struct {
	Wd                        string
	AppDir                    string
	Style                     string
	ModelIgnoreColumns        []string
	ModelMysqlDatasource      bool
	ModelMysqlDatasourceUrl   string
	ModelMysqlDatasourceTable []string
}

func (js *JzeroSql) Gen() error {
	dir := js.AppDir
	if dir == "" {
		dir = "."
	}

	if js.ModelMysqlDatasource {
		fmt.Printf("%s to generate model code from url %s.\n", color.WithColor("Start", color.FgGreen), js.ModelMysqlDatasourceUrl)
		// get tables from url
		tables := getMysqlAllTables(js.ModelMysqlDatasourceUrl)

		for _, table := range tables {
			fmt.Printf("%s table %s\n", color.WithColor("Using", color.FgGreen), table)
			command := fmt.Sprintf("goctl model mysql datasource --url '%s' --table %s --dir %s --home %s --style %s -i '%s'",
				js.ModelMysqlDatasourceUrl,
				table,
				filepath.Join(dir, "internal", "model", strings.ToLower(table)),
				filepath.Join(embeded.Home, "go-zero"),
				js.Style,
				strings.Join(js.ModelIgnoreColumns, ","))
			_, err := execx.Run(command, js.Wd)
			if err != nil {
				return err
			}
		}
		return nil
	}

	sqlDir := filepath.Join(js.Wd, js.AppDir, "desc", "sql")
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

				command := fmt.Sprintf("goctl model mysql ddl --src %s --dir %s --home %s --style %s -i '%s'",
					filepath.Join(dir, "desc", "sql", f.Name()),
					filepath.Join(dir, "internal", "model", strings.ToLower(f.Name()[0:len(f.Name())-len(path.Ext(f.Name()))])),
					filepath.Join(embeded.Home, "go-zero"),
					js.Style,
					strings.Join(js.ModelIgnoreColumns, ","))
				_, err = execx.Run(command, js.Wd)
				if err != nil {
					return err
				}
			}
		}
		fmt.Println(color.WithColor("Done", color.FgGreen))
	}
	return nil
}

type Table struct {
	Name string `db:"name"`
}

func getMysqlAllTables(url string) []string {
	sqlConn := sqlx.NewSqlConn("mysql", url)

	var tables []string
	err := sqlConn.QueryRowsCtx(context.Background(), &tables, "show tables")
	if err != nil {
		return nil
	}
	return tables
}
