package genmodel

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type ImportWithAlias struct {
	Alias string
	Path  string
}

type TableInfo struct {
	Alias             string
	Name              string
	FullName          string // 完整表名，如 "log.user"
	WithCache         bool
	HasCacheExpiry    bool
	HasNotFoundExpiry bool
}

func (jm *JzeroModel) GenRegister(tables []string) error {
	logx.Debugf("get register tables: %v", tables)

	slices.Sort(tables)

	var imports []string
	var importsWithAlias []ImportWithAlias
	var tablePackages []string
	var tableInfos []TableInfo

	mutiModels := make(map[string][]string)
	mutiModelsWithAlias := make(map[string][]TableInfo)

	for _, t := range tables {
		var isMutiModel bool
		tf := t
		if strings.Contains(t, ".") {
			tf = filepath.Join(strings.Split(t, ".")...)
			isMutiModel = true
		}
		mf := filepath.Join("internal", "model", strings.ToLower(tf))
		if !pathx.FileExists(mf) {
			logx.Debugf("%s table generated code not exists, skip", tf)
			continue
		}

		if isMutiModel {
			imports = append(imports, fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))))
			importsWithAlias = append(importsWithAlias, ImportWithAlias{
				Alias: strings.ReplaceAll(strings.ReplaceAll(t, ".", "_"), "-", "_"),
				Path:  fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))),
			})
			tableName := strings.ToLower(filepath.Base(tf))
			mutiModels[filepath.Dir(tf)] = append(mutiModels[filepath.Dir(tf)], tableName)
			// 将路径分隔符替换为点，如 "log/user" -> "log.user"
			fullName := strings.ReplaceAll(strings.ToLower(tf), "/", ".")
			mutiModelsWithAlias[filepath.Dir(tf)] = append(mutiModelsWithAlias[filepath.Dir(tf)], TableInfo{
				Alias:     strings.ReplaceAll(strings.ReplaceAll(t, ".", "_"), "-", "_"),
				Name:      tableName,
				FullName:  fullName, // 完整表名，如 "log.user"
				WithCache: getIsCacheTable(t),
			})
		} else {
			imports = append(imports, fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))))
			importsWithAlias = append(importsWithAlias, ImportWithAlias{
				Path: fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))),
			})
			tablePackages = append(tablePackages, strings.ToLower(tf))
			tableInfos = append(tableInfos, TableInfo{
				Name:      strings.ToLower(tf),
				WithCache: getIsCacheTable(t),
			})
		}
	}

	logx.Debugf("get register imports: %v", imports)
	logx.Debugf("get register table packages: %v", tablePackages)
	logx.Debugf("get register muti models: %v", mutiModels)

	// Build cache expiry table maps - only when cache is enabled
	var modelExpiryTable map[string]int64
	var modelNotFoundExpiryTable map[string]int64
	if config.C.Gen.ModelCache {
		modelExpiryTable = make(map[string]int64)
		modelNotFoundExpiryTable = make(map[string]int64)
		for _, v := range config.C.Gen.ModelCacheExpiryTable {
			if v.Expiry > 0 {
				modelExpiryTable[v.Table] = v.Expiry
			}
			if v.NotFoundExpiry > 0 {
				modelNotFoundExpiryTable[v.Table] = v.NotFoundExpiry
			}
		}
	}

	// Update tableInfos with cache expiry info
	for i := range tableInfos {
		if modelExpiryTable != nil {
			if _, ok := modelExpiryTable[tableInfos[i].Name]; ok {
				tableInfos[i].HasCacheExpiry = true
			}
		}
		if modelNotFoundExpiryTable != nil {
			if _, ok := modelNotFoundExpiryTable[tableInfos[i].Name]; ok {
				tableInfos[i].HasNotFoundExpiry = true
			}
		}
	}

	// Update mutiModelsWithAlias with cache expiry info
	for k := range mutiModelsWithAlias {
		for i := range mutiModelsWithAlias[k] {
			fullName := mutiModelsWithAlias[k][i].FullName
			if modelExpiryTable != nil {
				if _, ok := modelExpiryTable[fullName]; ok {
					mutiModelsWithAlias[k][i].HasCacheExpiry = true
				}
			}
			if modelNotFoundExpiryTable != nil {
				if _, ok := modelNotFoundExpiryTable[fullName]; ok {
					mutiModelsWithAlias[k][i].HasNotFoundExpiry = true
				}
			}
		}
	}

	template, err := templatex.ParseTemplate(filepath.Join("model", "model.go.tpl"), map[string]any{
		"Imports":                  imports,
		"ImportsWithAlias":         importsWithAlias,
		"TablePackages":            tablePackages,
		"TableInfos":               tableInfos,
		"MutiModels":               mutiModels,          // 兼容
		"MutiModelsWithAlias":      mutiModelsWithAlias, // 兼容
		"ModelExpiryTable":         modelExpiryTable,
		"ModelNotFoundExpiryTable": modelNotFoundExpiryTable,
		"ModelCache":               config.C.Gen.ModelCache,
		"ModelNewOriginal":         config.C.Gen.ModelNewOriginal,
	}, lo.If(
		// 兼容老版本 model 路径
		// TODO: wait to remove
		embeded.ReadTemplateFile(filepath.Join("plugins", "model", "model.go.tpl")) != nil,
		embeded.ReadTemplateFile(filepath.Join("plugins", "model", "model.go.tpl"))).
		Else(
			embeded.ReadTemplateFile(filepath.Join("model", "model.go.tpl")),
		),
	)
	if err != nil {
		return err
	}

	format, err := gosimports.Process("", template, nil)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("internal", "model", "model.go"), format, 0o644)
}
