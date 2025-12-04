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

func (jm *JzeroModel) GenRegister(tables []string) error {
	logx.Debugf("get register tables: %v", tables)

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
		mf := filepath.Join("internal", "model", strings.ToLower(t))
		if !pathx.FileExists(mf) {
			logx.Debugf("%s table generated code not exists, skip", t)
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

	logx.Debugf("get register imports: %v", imports)
	logx.Debugf("get register table packages: %v", tablePackages)
	logx.Debugf("get register muti models: %v", mutiModels)

	template, err := templatex.ParseTemplate(filepath.Join("model", "model.go.tpl"), map[string]any{
		"Imports":       imports,
		"TablePackages": tablePackages,
		"MutiModels":    mutiModels,
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
