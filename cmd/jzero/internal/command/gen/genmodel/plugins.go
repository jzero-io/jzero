package genmodel

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func (jm *JzeroModel) GenRegister(tables []string) error {
	slices.Sort(tables)

	var imports []string
	var tablePackages []string

	for _, t := range tables {
		if !pathx.FileExists(fmt.Sprintf("internal/model/%s", strings.ToLower(t))) {
			continue
		}
		imports = append(imports, fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(t)))
		tablePackages = append(tablePackages, strings.ToLower(t))
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Imports":       imports,
		"TablePackages": tablePackages,
		"withCache":     config.C.Gen.ModelCache,
	}, embeded.ReadTemplateFile(filepath.Join("plugins", "model", "model.go.tpl")))
	if err != nil {
		return err
	}

	format, err := gosimports.Process("", template, &gosimports.Options{
		Comments:   true,
		FormatOnly: true,
	})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("internal", "model", "model.go"), format, 0o644)
}

func (jm *JzeroModel) GenDDL(sqlConn sqlx.SqlConn, tables []string) ([]string, error) {
	var (
		eg          errgroup.Group
		tableDDLMap sync.Map
	)
	// TODO: add flag to manage concurrency
	eg.SetLimit(1)
	for _, t := range tables {
		ct := t
		eg.Go(func() error {
			ddl, err := getTableDDL(sqlConn, config.C.Gen.ModelDriver, ct)
			if err != nil {
				return err
			}
			re := regexp.MustCompile(`AUTO_INCREMENT=\d+\s*`)
			ddl = re.ReplaceAllString(ddl, "")
			tableDDLMap.Store(ct, ddl)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Join("desc", "sql"), 0o755); err != nil {
		return nil, err
	}

	var writeTables []string
	for _, v := range tables {
		if s, ok := tableDDLMap.Load(v); ok {
			if len(config.C.Gen.ModelDatasourceTable) != 0 && config.C.Gen.ModelDatasourceTable[0] != "*" {
				if lo.Contains(config.C.Gen.ModelDatasourceTable, cast.ToString(v)) {
					writeTables = append(writeTables, filepath.Join("desc", "sql", fmt.Sprintf("%s.sql", v)))
					if err := os.WriteFile(filepath.Join("desc", "sql", fmt.Sprintf("%s.sql", v)), []byte(cast.ToString(s)), 0o644); err != nil {
						return nil, err
					}
				}
			} else if len(config.C.Gen.ModelDatasourceTable) != 0 && config.C.Gen.ModelDatasourceTable[0] == "*" {
				writeTables = append(writeTables, filepath.Join("desc", "sql", fmt.Sprintf("%s.sql", v)))
				if err := os.WriteFile(filepath.Join("desc", "sql", fmt.Sprintf("%s.sql", v)), []byte(cast.ToString(s)), 0o644); err != nil {
					return nil, err
				}
			}
		}
	}

	return writeTables, nil
}
