package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

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

	jzeroRpc := JzeroRpc{Wd: wd, Module: moduleStruct.Path, Style: Style, RemoveSuffix: RemoveSuffix}
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
	_ = os.Remove(filepath.Join(wd, fmt.Sprintf("%s.go", GetApiServiceName(filepath.Join(wd, "desc", "api")))))
	_ = os.Remove(filepath.Join(wd, "etc", fmt.Sprintf("%s.yaml", GetApiServiceName(filepath.Join(wd, "desc", "api")))))
	protoFilenames, err := GetProtoFilepath(filepath.Join(wd, "desc", "proto"))
	if err == nil {
		for _, v := range protoFilenames {
			v = filepath.Base(v)
			fileBase := v[0 : len(v)-len(path.Ext(v))]
			rmf := strings.ReplaceAll(strings.ToLower(fileBase), "-", "")
			rmf = strings.ReplaceAll(rmf, "_", "")
			_ = os.Remove(filepath.Join(wd, fmt.Sprintf("%s.go", rmf)))
			_ = os.Remove(filepath.Join(wd, "etc", fmt.Sprintf("%s.yaml", rmf)))
		}
	}
}

func init() {
	logx.Disable()
}
