package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/pkg/errors"

	"github.com/jzero-io/jzero/pkg/mod"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

var (
	Style              string
	RemoveSuffix       bool
	ChangeReplaceTypes bool
)

var (
	// ModelMysqlIgnoreColumns goctl model flags
	ModelMysqlIgnoreColumns []string

	ModelMysqlCache       bool
	ModelMysqlCachePrefix string

	// ModelMysqlDatasource goctl model datasource
	ModelMysqlDatasource      bool
	ModelMysqlDatasourceUrl   string
	ModelMysqlDatasourceTable []string
)

type ApiFileTypes struct {
	Filepath string
	ApiSpec  spec.ApiSpec
	GenTypes []spec.Type

	Base bool
}

func Gen(_ *cobra.Command, _ []string) error {
	wd, err := os.Getwd()
	cobra.CheckErr(err)
	fmt.Printf("%s working dir %s\n", color.WithColor("Enter", color.FgGreen), wd)

	moduleStruct, err := mod.GetGoMod(wd)
	cobra.CheckErr(errors.Wrapf(err, "get go module struct error"))

	defer func() {
		RemoveExtraFiles(wd)
	}()

	jzeroRpc := JzeroRpc{
		Wd:                 wd,
		Module:             moduleStruct.Path,
		Style:              Style,
		RemoveSuffix:       RemoveSuffix,
		ChangeReplaceTypes: ChangeReplaceTypes,
		Etc:                filepath.Join("etc", "etc.yaml"),
	}
	err = jzeroRpc.Gen()
	if err != nil {
		return err
	}

	jzeroApi := JzeroApi{Wd: wd, Module: moduleStruct.Path, Style: Style, RemoveSuffix: RemoveSuffix, ChangeReplaceTypes: ChangeReplaceTypes}
	err = jzeroApi.Gen()
	if err != nil {
		return err
	}

	jzeroSql := JzeroSql{
		Wd:                        wd,
		Style:                     Style,
		ModelIgnoreColumns:        ModelMysqlIgnoreColumns,
		ModelMysqlDatasource:      ModelMysqlDatasource,
		ModelMysqlDatasourceUrl:   ModelMysqlDatasourceUrl,
		ModelMysqlDatasourceTable: ModelMysqlDatasourceTable,
		ModelMysqlCache:           ModelMysqlCache,
		ModelMysqlCachePrefix:     ModelMysqlCachePrefix,
	}
	err = jzeroSql.Gen()
	if err != nil {
		return err
	}

	return nil
}

func RemoveExtraFiles(wd string) {
	_ = os.Remove(filepath.Join(wd, getApiFrameMainGoFilename(wd)))
	_ = os.Remove(filepath.Join(wd, "etc", getApiFrameEtcFilename(wd)))

	protoFilenames, err := GetProtoFilepath(filepath.Join("desc", "proto"))
	if err == nil {
		for _, v := range protoFilenames {
			v = filepath.Base(v)
			fileBase := v[0 : len(v)-len(path.Ext(v))]
			_ = os.Remove(filepath.Join(wd, getProtoFrameMainGoFilename(fileBase)))
			_ = os.Remove(filepath.Join(wd, "etc", getProtoFrameEtcFilename(fileBase)))
		}
	}
}

// getApiFrameMainGoFilename: goctl/api/gogen/genmain.go
func getApiFrameMainGoFilename(wd string) string {
	serviceName := GetApiServiceName(filepath.Join(wd, "desc", "api"))
	serviceName = strings.ToLower(serviceName)
	filename, err := format.FileNamingFormat(Style, serviceName)
	if err != nil {
		return ""
	}

	if strings.HasSuffix(filename, "-api") {
		filename = strings.ReplaceAll(filename, "-api", "")
	}
	return filename + ".go"
}

// getApiFrameEtcFilename: goctl/api/gogen/genetc.go
func getApiFrameEtcFilename(wd string) string {
	serviceName := GetApiServiceName(filepath.Join(wd, "desc", "api"))
	filename, err := format.FileNamingFormat(Style, serviceName)
	if err != nil {
		return ""
	}
	return filename + ".yaml"
}

// getProtoFrameMainGoFilename: goctl/rpc/generator/genmain.go
func getProtoFrameMainGoFilename(source string) string {
	filename, err := format.FileNamingFormat(Style, source)
	if err != nil {
		return ""
	}
	return filename + ".go"
}

// getProtoFrameEtcFilename: goctl/rpc/generator/genetc.go
func getProtoFrameEtcFilename(source string) string {
	filename, err := format.FileNamingFormat(Style, source)
	if err != nil {
		return ""
	}
	return filename + ".yaml"
}

func init() {
	logx.Disable()
}
