package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen/genapi"
	"github.com/jzero-io/jzero/internal/gen/genmodel"
	"github.com/jzero-io/jzero/internal/gen/genrpc"
	"github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/mod"
)

func Gen(c config.Config) error {
	fmt.Printf("%s working dir %s\n", color.WithColor("Enter", color.FgGreen), c.Wd())

	var module string
	moduleStruct, err := mod.GetGoMod(c.Wd())
	if err != nil {
		return errors.Wrapf(err, "get go module struct error")
	}
	module = moduleStruct.Path

	if !pathx.FileExists("go.mod") {
		module, err = mod.GetParentPackage(c.Wd())
		if err != nil {
			return errors.Wrapf(err, "get parent package error")
		}
	}

	defer func() {
		RemoveExtraFiles(c.Wd(), c.Gen.Style)
	}()

	jzeroRpc := genrpc.JzeroRpc{
		Wd:               c.Wd(),
		Module:           module,
		Style:            c.Gen.Style,
		RemoveSuffix:     c.Gen.RemoveSuffix,
		ChangeLogicTypes: c.Gen.ChangeLogicTypes,
		RpcStylePatch:    c.Gen.RpcStylePatch,
	}
	err = jzeroRpc.Gen()
	if err != nil {
		return err
	}

	jzeroApi := genapi.JzeroApi{
		Wd:                 c.Wd(),
		Module:             module,
		Style:              c.Gen.Style,
		RemoveSuffix:       c.Gen.RemoveSuffix,
		ChangeReplaceTypes: c.Gen.ChangeLogicTypes,
		RegenApiHandler:    c.Gen.RegenApiHandler,
		SplitApiTypesDir:   c.Gen.SplitApiTypesDir,
		GitDiff:            c.Gen.GitDiff,
		ApiGitDiffPath:     c.Gen.ApiGitDiffPath,
	}
	err = jzeroApi.Gen()
	if err != nil {
		return err
	}

	jzeroSql := genmodel.JzeroSql{
		Wd:                        c.Wd(),
		Style:                     c.Gen.Style,
		ModelStrict:               c.Gen.ModelMysqlStrict,
		ModelIgnoreColumns:        c.Gen.ModelMysqlIgnoreColumns,
		ModelMysqlDatasource:      c.Gen.ModelMysqlDatasource,
		ModelMysqlDatasourceUrl:   c.Gen.ModelMysqlDatasourceUrl,
		ModelMysqlDatasourceTable: c.Gen.ModelMysqlDatasourceTable,
		ModelMysqlCache:           c.Gen.ModelMysqlCache,
		ModelMysqlCachePrefix:     c.Gen.ModelMysqlCachePrefix,
		GitDiff:                   c.Gen.GitDiff,
		ModelGitDiffPath:          c.Gen.ModelGitDiffPath,
	}
	err = jzeroSql.Gen()
	if err != nil {
		return err
	}

	return nil
}

func RemoveExtraFiles(wd, style string) {
	if pathx.FileExists(filepath.Join("desc", "api")) {
		if desc.GetApiFrameMainGoFilename(wd, style) != "main.go" {
			if err := os.Remove(filepath.Join(wd, desc.GetApiFrameMainGoFilename(wd, style))); err != nil {
				logx.Debugf("remove api frame main go file error: %s", err.Error())
			}
		}
		if desc.GetApiFrameEtcFilename(wd, style) != "etc.yaml" {
			if err := os.Remove(filepath.Join(wd, "etc", desc.GetApiFrameEtcFilename(wd, style))); err != nil {
				logx.Debugf("remove api etc file error: %s", err.Error())
			}
		}
	}

	if pathx.FileExists(filepath.Join("desc", "proto")) {
		protoFilenames, err := desc.GetProtoFilepath(filepath.Join("desc", "proto"))
		if err == nil {
			for _, v := range protoFilenames {
				v = filepath.Base(v)
				fileBase := v[0 : len(v)-len(path.Ext(v))]
				if desc.GetProtoFrameMainGoFilename(fileBase, style) != "main.go" {
					if err = os.Remove(filepath.Join(wd, desc.GetProtoFrameMainGoFilename(fileBase, style))); err != nil {
						logx.Debugf("remove proto frame main go file error: %s", err.Error())
					}
				}
				if desc.GetProtoFrameEtcFilename(fileBase, style) != "etc.yaml" {
					if err = os.Remove(filepath.Join(wd, "etc", desc.GetProtoFrameEtcFilename(fileBase, style))); err != nil {
						logx.Debugf("remove proto etc file error: %s", err.Error())
					}
				}
			}
		}
	}
}
