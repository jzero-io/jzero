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

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

func (jm *JzeroMongo) GenRegister(types []string) error {
	logx.Debugf("get register tables: %v", types)

	slices.Sort(types)

	var imports []string
	var typePackages []string

	mutiModels := make(map[string][]string)

	for _, t := range types {
		var isMutiModel bool
		if strings.Contains(t, ".") {
			t = filepath.Join(strings.Split(t, ".")...)
			isMutiModel = true
		}
		mf := filepath.Join("internal", "mongo", strings.ToLower(t))
		if !pathx.FileExists(mf) {
			logx.Debugf("%s mongo model generated code not exists, skip", t)
			continue
		}

		imports = append(imports, fmt.Sprintf("%s/internal/mongo/%s", jm.Module, strings.ToLower(filepath.ToSlash(t))))

		if isMutiModel {
			mutiModels[filepath.Dir(t)] = append(mutiModels[filepath.Dir(t)], strings.ToLower(filepath.Base(t)))
		} else {
			typePackages = append(typePackages, strings.ToLower(t))
		}
	}

	logx.Debugf("get register imports: %v", imports)
	logx.Debugf("get register types packages: %v", typePackages)
	logx.Debugf("get register muti models: %v", mutiModels)

	// Determine if any mongo types have cache enabled
	mongoHasCache := config.C.Gen.MongoCache
	if len(config.C.Gen.MongoCacheType) > 0 {
		mongoHasCache = len(config.C.Gen.MongoCacheType) > 0
	}

	template, err := templatex.ParseTemplate(filepath.Join("plugins", "mongo", "model.go.tpl"), map[string]any{
		"Imports":      imports,
		"TypePackages": typePackages,
		"MutiModels":   mutiModels,
		"withCache":    mongoHasCache,
	}, embeded.ReadTemplateFile(filepath.Join("plugins", "mongo", "model.go.tpl")))
	if err != nil {
		return err
	}

	format, err := gosimports.Process("", template, nil)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("internal", "mongo", "model.go"), format, 0o644)
}
