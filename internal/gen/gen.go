package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/mod"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

var (
	WorkingDir string
	AppDir     string

	Version            string
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
	// change dir
	if WorkingDir != "" {
		err := os.Chdir(WorkingDir)
		cobra.CheckErr(err)
	}

	wd, err := os.Getwd()
	cobra.CheckErr(err)
	fmt.Printf("%s working dir %s\n", color.WithColor("Enter", color.FgGreen), wd)

	if embeded.Home == "" {
		home, _ := os.UserHomeDir()
		embeded.Home = filepath.Join(home, ".jzero", Version)
	}

	moduleStruct, err := mod.GetGoMod(wd)
	cobra.CheckErr(errors.Wrapf(err, "get go module struct error"))

	defer func() {
		removeExtraFiles(wd, AppDir)
	}()

	jzeroRpc := JzeroRpc{Wd: wd, AppDir: AppDir, Module: moduleStruct.Path, Style: Style, RemoveSuffix: RemoveSuffix}
	err = jzeroRpc.Gen()
	if err != nil {
		return err
	}

	jzeroApi := JzeroApi{Wd: wd, AppDir: AppDir, Module: moduleStruct.Path, Style: Style, RemoveSuffix: RemoveSuffix, ChangeReplaceTypes: ChangeReplaceTypes}
	err = jzeroApi.Gen()
	if err != nil {
		return err
	}

	jzeroSql := JzeroSql{
		Wd:                        wd,
		AppDir:                    AppDir,
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

func removeExtraFiles(wd string, appDir string) {
	_ = os.Remove(filepath.Join(wd, appDir, fmt.Sprintf("%s.go", GetApiServiceName(filepath.Join(wd, appDir, "desc", "api")))))
	_ = os.Remove(filepath.Join(wd, appDir, "etc", fmt.Sprintf("%s.yaml", GetApiServiceName(filepath.Join(wd, appDir, "desc", "api")))))
	protoFilenames, err := GetProtoFilenames(filepath.Join(wd, appDir, "desc", "proto"))
	if err == nil {
		for _, v := range protoFilenames {
			fileBase := v[0 : len(v)-len(path.Ext(v))]
			rmf := strings.ReplaceAll(strings.ToLower(fileBase), "-", "")
			rmf = strings.ReplaceAll(rmf, "_", "")
			_ = os.Remove(filepath.Join(wd, appDir, fmt.Sprintf("%s.go", rmf)))
			_ = os.Remove(filepath.Join(wd, appDir, "etc", fmt.Sprintf("%s.yaml", rmf)))
		}
	}
}

func init() {
	logx.Disable()
}
