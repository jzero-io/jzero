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

	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type ImportWithAlias struct {
	Alias string
	Path  string
}

type NameWithAlias struct {
	Alias string
	Name  string
}

func (jm *JzeroModel) GenRegister(tables []string) error {
	logx.Debugf("get register tables: %v", tables)

	slices.Sort(tables)

	var imports []string
	var importsWithAlias []ImportWithAlias
	var tablePackages []string

	mutiModels := make(map[string][]string)
	mutiModelsWithAlias := make(map[string][]NameWithAlias)

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
			mutiModels[filepath.Dir(tf)] = append(mutiModels[filepath.Dir(tf)], strings.ToLower(filepath.Base(tf)))
			mutiModelsWithAlias[filepath.Dir(tf)] = append(mutiModelsWithAlias[filepath.Dir(tf)], NameWithAlias{
				Alias: strings.ReplaceAll(strings.ReplaceAll(t, ".", "_"), "-", "_"),
				Name:  strings.ToLower(filepath.Base(tf)),
			})
		} else {
			imports = append(imports, fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))))
			importsWithAlias = append(importsWithAlias, ImportWithAlias{
				Path: fmt.Sprintf("%s/internal/model/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))),
			})
			tablePackages = append(tablePackages, strings.ToLower(tf))
		}
	}

	logx.Debugf("get register imports: %v", imports)
	logx.Debugf("get register table packages: %v", tablePackages)
	logx.Debugf("get register muti models: %v", mutiModels)

	template, err := templatex.ParseTemplate(filepath.Join("model", "model.go.tpl"), map[string]any{
		"Imports":             imports,
		"ImportsWithAlias":    importsWithAlias,
		"TablePackages":       tablePackages,
		"MutiModels":          mutiModels,          // 兼容
		"MutiModelsWithAlias": mutiModelsWithAlias, // 兼容
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
