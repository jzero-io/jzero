package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/pkg/mod"
)

func Gen(gcf config.GenConfig) error {
	fmt.Printf("%s working dir %s\n", color.WithColor("Enter", color.FgGreen), gcf.Wd())

	var module string
	moduleStruct, err := mod.GetGoMod(gcf.Wd())
	if err != nil {
		return errors.Wrapf(err, "get go module struct error")
	}
	module = moduleStruct.Path

	if !pathx.FileExists("go.mod") {
		module = filepath.ToSlash(filepath.Join(module, filepath.Base(gcf.Wd())))
	}

	defer func() {
		RemoveExtraFiles(gcf.Wd(), gcf.Style)
	}()

	jzeroRpc := JzeroRpc{
		Wd:               gcf.Wd(),
		Module:           module,
		Style:            gcf.Style,
		RemoveSuffix:     gcf.RemoveSuffix,
		ChangeLogicTypes: gcf.ChangeLogicTypes,
		RpcStylePatch:    gcf.RpcStylePatch,
	}
	err = jzeroRpc.Gen()
	if err != nil {
		return err
	}

	jzeroApi := JzeroApi{
		Wd:                 gcf.Wd(),
		Module:             module,
		Style:              gcf.Style,
		RemoveSuffix:       gcf.RemoveSuffix,
		ChangeReplaceTypes: gcf.ChangeLogicTypes,
		RegenApiHandler:    gcf.RegenApiHandler,
		RegenApiTypes:      gcf.RegenApiTypes,
		SplitApiTypesDir:   gcf.SplitApiTypesDir,
	}
	err = jzeroApi.Gen()
	if err != nil {
		return err
	}

	jzeroSql := JzeroSql{
		Wd:                        gcf.Wd(),
		Style:                     gcf.Style,
		ModelStrict:               gcf.ModelMysqlStrict,
		ModelIgnoreColumns:        gcf.ModelMysqlIgnoreColumns,
		ModelMysqlDatasource:      gcf.ModelMysqlDatasource,
		ModelMysqlDatasourceUrl:   gcf.ModelMysqlDatasourceUrl,
		ModelMysqlDatasourceTable: gcf.ModelMysqlDatasourceTable,
		ModelMysqlCache:           gcf.ModelMysqlCache,
		ModelMysqlCachePrefix:     gcf.ModelMysqlCachePrefix,
		ModelGitDiff:              gcf.ModelGitDiff,
	}
	err = jzeroSql.Gen()
	if err != nil {
		return err
	}

	return nil
}

func RemoveExtraFiles(wd, style string) {
	if pathx.FileExists(filepath.Join("desc", "api")) {
		if err := os.Remove(filepath.Join(wd, getApiFrameMainGoFilename(wd, style))); err != nil {
			logx.Debugf("remove api frame main go file error: %s", err.Error())
		}
		if err := os.Remove(filepath.Join(wd, "etc", getApiFrameEtcFilename(wd, style))); err != nil {
			logx.Debugf("remove api etc file error: %s", err.Error())
		}
	}

	if pathx.FileExists(filepath.Join("desc", "proto")) {
		protoFilenames, err := GetProtoFilepath(filepath.Join("desc", "proto"))
		if err == nil {
			for _, v := range protoFilenames {
				v = filepath.Base(v)
				fileBase := v[0 : len(v)-len(path.Ext(v))]
				if err = os.Remove(filepath.Join(wd, getProtoFrameMainGoFilename(fileBase, style))); err != nil {
					logx.Debugf("remove proto frame main go file error: %s", err.Error())
				}
				if err = os.Remove(filepath.Join(wd, "etc", getProtoFrameEtcFilename(fileBase, style))); err != nil {
					logx.Debugf("remove proto etc file error: %s", err.Error())
				}
			}
		}
	}
}

// getApiFrameMainGoFilename: goctl/api/gogen/genmain.go
func getApiFrameMainGoFilename(wd, style string) string {
	serviceName := GetApiServiceName(filepath.Join(wd, "desc", "api"))
	serviceName = strings.ToLower(serviceName)
	filename, err := format.FileNamingFormat(style, serviceName)
	if err != nil {
		return ""
	}

	if strings.HasSuffix(filename, "-api") {
		filename = strings.ReplaceAll(filename, "-api", "")
	}
	return filename + ".go"
}

// getApiFrameEtcFilename: goctl/api/gogen/genetc.go
func getApiFrameEtcFilename(wd, style string) string {
	serviceName := GetApiServiceName(filepath.Join(wd, "desc", "api"))
	filename, err := format.FileNamingFormat(style, serviceName)
	if err != nil {
		return ""
	}
	return filename + ".yaml"
}

// getProtoFrameMainGoFilename: goctl/rpc/generator/genmain.go
func getProtoFrameMainGoFilename(source, style string) string {
	filename, err := format.FileNamingFormat(style, source)
	if err != nil {
		return ""
	}
	return filename + ".go"
}

// getProtoFrameEtcFilename: goctl/rpc/generator/genetc.go
func getProtoFrameEtcFilename(source, style string) string {
	filename, err := format.FileNamingFormat(style, source)
	if err != nil {
		return ""
	}
	return filename + ".yaml"
}
