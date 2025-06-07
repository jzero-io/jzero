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

	mutiModels := make(map[string][]string)

	for _, t := range tables {
		var isMutiModel bool
		if strings.Contains(t, ".") {
			t = filepath.Join(strings.Split(t, ".")...)
			isMutiModel = true
		}
		mf := filepath.Join("internal", "model", t)
		if !pathx.FileExists(mf) {
			continue
		}

		imports = append(imports, fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(filepath.ToSlash(t))))

		if isMutiModel {
			imports = append(imports, fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(filepath.ToSlash(t))))
			mutiModels[filepath.Dir(t)] = append(mutiModels[filepath.Dir(t)], strings.ToLower(filepath.Base(t)))
		} else {
			tablePackages = append(tablePackages, strings.ToLower(t))
		}
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Imports":       imports,
		"TablePackages": tablePackages,
		"MutiModels":    mutiModels,
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

func (jm *JzeroModel) GenDDL(conns []Conn, tables []string) ([]string, error) {
	var (
		eg          errgroup.Group
		tableDDLMap sync.Map
	)
	eg.SetLimit(len(tables))

	for _, conn := range conns {
		eg.Go(func() error {
			for _, t := range tables {
				eg.Go(func() error {
					ddl, err := getTableDDL(conn.SqlConn, config.C.Gen.ModelDriver, t)
					if err != nil {
						return err
					}
					re := regexp.MustCompile(`AUTO_INCREMENT=\d+\s*`)
					ddl = re.ReplaceAllString(ddl, "")
					tableDDLMap.Store(t, ddl)
					return nil
				})
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Join("desc", "sql"), 0o755); err != nil {
		return nil, err
	}

	var (
		writeTables []string
	)

	if len(config.C.Gen.ModelDatasourceTable) == 1 && config.C.Gen.ModelDatasourceTable[0] == "*" {
		config.C.Gen.ModelDatasourceTable = tables
		if len(config.C.Gen.ModelDatasourceUrl) == 1 {
			config.C.Gen.ModelDatasourceTable = lo.Map(tables, func(item string, index int) string {
				return strings.Split(item, ".")[1]
			})
		}
	}

	for _, v := range config.C.Gen.ModelDatasourceTable {
		var vWithScheme string
		if strings.Contains(v, ".") {
			vWithScheme = v
		} else {
			vWithScheme = fmt.Sprintf("%s.%s", conns[0].Scheme, v)
		}
		if s, ok := tableDDLMap.Load(vWithScheme); ok {
			if len(config.C.Gen.ModelDatasourceTable) != 0 && config.C.Gen.ModelDatasourceTable[0] != "*" {
				if lo.Contains(tables, vWithScheme) {
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
