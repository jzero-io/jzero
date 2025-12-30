package genmongo

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rinchsan/gosimports"
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

func (jm *JzeroMongo) GenRegister(types []string) error {
	logx.Debugf("get register tables: %v", types)

	slices.Sort(types)

	var imports []string
	var importsWithAlias []ImportWithAlias
	var typePackages []string

	mutiModels := make(map[string][]string)
	mutiModelsWithAlias := make(map[string][]NameWithAlias)

	for _, t := range types {
		var isMutiModel bool
		tf := t
		if strings.Contains(t, ".") {
			tf = filepath.Join(strings.Split(t, ".")...)
			isMutiModel = true
		}
		mf := filepath.Join("internal", "mongo", strings.ToLower(tf))
		if !pathx.FileExists(mf) {
			logx.Debugf("%s mongo model generated code not exists, skip", t)
			continue
		}

		if isMutiModel {
			imports = append(imports, fmt.Sprintf("%s/internal/mongo/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))))
			importsWithAlias = append(importsWithAlias, ImportWithAlias{
				Alias: strings.ReplaceAll(strings.ReplaceAll(t, ".", "_"), "-", "_"),
				Path:  fmt.Sprintf("%s/internal/mongo/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))),
			})
			mutiModels[filepath.Dir(tf)] = append(mutiModels[filepath.Dir(tf)], strings.ToLower(filepath.Base(tf)))
			mutiModelsWithAlias[filepath.Dir(tf)] = append(mutiModelsWithAlias[filepath.Dir(tf)], NameWithAlias{
				Alias: strings.ReplaceAll(strings.ReplaceAll(t, ".", "_"), "-", "_"),
				Name:  strings.ToLower(filepath.Base(tf)),
			})
		} else {
			imports = append(imports, fmt.Sprintf("%s/internal/mongo/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))))
			importsWithAlias = append(importsWithAlias, ImportWithAlias{
				Alias: filepath.Base(tf),
				Path:  fmt.Sprintf("%s/internal/mongo/%s", jm.Module, strings.ToLower(filepath.ToSlash(tf))),
			})
			typePackages = append(typePackages, strings.ToLower(tf))
		}
	}

	logx.Debugf("get register imports: %v", imports)
	logx.Debugf("get register types packages: %v", typePackages)
	logx.Debugf("get register muti models: %v", mutiModels)

	template, err := templatex.ParseTemplate(filepath.Join("mongo", "model.go.tpl"), map[string]any{
		"Imports":             imports,
		"ImportsWithAlias":    importsWithAlias,
		"TypePackages":        typePackages,
		"MutiModels":          mutiModels,          // 兼容
		"MutiModelsWithAlias": mutiModelsWithAlias, // 兼容
	}, embeded.ReadTemplateFile(filepath.Join("mongo", "model.go.tpl")))
	if err != nil {
		return err
	}

	format, err := gosimports.Process("", template, nil)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("internal", "mongo", "model.go"), format, 0o644)
}
