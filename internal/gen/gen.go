package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/pkg/mod"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

func Gen(gcf config.GenConfig) error {
	fmt.Printf("%s working dir %s\n", color.WithColor("Enter", color.FgGreen), gcf.Wd())

	moduleStruct, err := mod.GetGoMod(gcf.Wd())
	if err != nil {
		return errors.Wrapf(err, "get go module struct error")
	}

	defer func() {
		RemoveExtraFiles(gcf.Wd(), gcf.Style)
	}()

	jzeroRpc := JzeroRpc{
		Wd:                 gcf.Wd(),
		Module:             moduleStruct.Path,
		Style:              gcf.Style,
		RemoveSuffix:       gcf.RemoveSuffix,
		ChangeReplaceTypes: gcf.ChangeReplaceTypes,
	}
	err = jzeroRpc.Gen()
	if err != nil {
		return err
	}

	jzeroApi := JzeroApi{
		Wd:                 gcf.Wd(),
		Module:             moduleStruct.Path,
		Style:              gcf.Style,
		RemoveSuffix:       gcf.RemoveSuffix,
		ChangeReplaceTypes: gcf.ChangeReplaceTypes,
	}
	err = jzeroApi.Gen()
	if err != nil {
		return err
	}

	jzeroSql := JzeroSql{
		Wd:                        gcf.Wd(),
		Style:                     gcf.Style,
		ModelIgnoreColumns:        gcf.ModelMysqlIgnoreColumns,
		ModelMysqlDatasource:      gcf.ModelMysqlDatasource,
		ModelMysqlDatasourceUrl:   gcf.ModelMysqlDatasourceUrl,
		ModelMysqlDatasourceTable: gcf.ModelMysqlDatasourceTable,
		ModelMysqlCache:           gcf.ModelMysqlCache,
		ModelMysqlCachePrefix:     gcf.ModelMysqlCachePrefix,
	}
	err = jzeroSql.Gen()
	if err != nil {
		return err
	}

	return nil
}

func RemoveExtraFiles(wd string, style string) {
	if err := os.Remove(filepath.Join(wd, getApiFrameMainGoFilename(wd, style))); err != nil {
		logx.Debugf("remove api frame main go file error: %s", err.Error())
	}
	if err := os.Remove(filepath.Join(wd, "etc", getApiFrameEtcFilename(wd, style))); err != nil {
		logx.Debugf("remove api etc file error: %s", err.Error())
	}

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

// getApiFrameMainGoFilename: goctl/api/gogen/genmain.go
func getApiFrameMainGoFilename(wd string, style string) string {
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
func getApiFrameEtcFilename(wd string, style string) string {
	serviceName := GetApiServiceName(filepath.Join(wd, "desc", "api"))
	filename, err := format.FileNamingFormat(style, serviceName)
	if err != nil {
		return ""
	}
	return filename + ".yaml"
}

// getProtoFrameMainGoFilename: goctl/rpc/generator/genmain.go
func getProtoFrameMainGoFilename(source string, style string) string {
	filename, err := format.FileNamingFormat(style, source)
	if err != nil {
		return ""
	}
	return filename + ".go"
}

// getProtoFrameEtcFilename: goctl/rpc/generator/genetc.go
func getProtoFrameEtcFilename(source string, style string) string {
	filename, err := format.FileNamingFormat(style, source)
	if err != nil {
		return ""
	}
	return filename + ".yaml"
}
