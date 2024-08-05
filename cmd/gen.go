/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/internal/gen/gendocs"
	"github.com/jzero-io/jzero/internal/gen/gensdk"
	"github.com/jzero-io/jzero/internal/gen/genswagger"
	"github.com/jzero-io/jzero/internal/gen/genzrpcclient"
	"github.com/jzero-io/jzero/pkg/mod"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
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

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: `Used to generate server/client code`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.SetConfig(CfgFile, cmd.Use, cmd.Flags())
		if err != nil {
			return err
		}

		// check go-zero api template
		home, _ := os.UserHomeDir()
		if !pathx.FileExists(filepath.Join(home, ".jzero", Version, "go-zero")) {
			err := embeded.WriteTemplateDir(filepath.Join("go-zero"), filepath.Join(home, ".jzero", Version, "go-zero"))
			cobra.CheckErr(err)
		}

		if !pathx.FileExists(c.Gen.Home) {
			home, _ := os.UserHomeDir()
			embeded.Home = filepath.Join(home, ".jzero", Version)
		}
		return gen.Gen(c.Gen)
	},
	SilenceUsage: true,
}

// genZRpcClientCmd represents the rpcClient command
var genZRpcClientCmd = &cobra.Command{
	Use:   "zrpcclient",
	Short: `Gen zrpc client code by proto`,
	PreRun: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		cobra.CheckErr(err)
		mod, err := mod.GetGoMod(wd)
		cobra.CheckErr(err)
		if genzrpcclient.Scope == "" {
			genzrpcclient.Scope = filepath.Base(mod.Path)
			genzrpcclient.Scope = strings.ReplaceAll(genzrpcclient.Scope, "-", "_")
		}

		if genzrpcclient.GoModule == "" {
			genzrpcclient.GoModule = filepath.ToSlash(filepath.Join(mod.Path, genzrpcclient.Output))
		} else {
			genzrpcclient.GenModule = true
		}

		if genzrpcclient.GoPackage == "" {
			genzrpcclient.GoPackage = strings.ReplaceAll(strings.ToLower(filepath.Base(genzrpcclient.GoModule)), "-", "_")
		}
	},
	RunE: genzrpcclient.Generate,
}

// genSwaggerCmd represents the genSwagger command
var genSwaggerCmd = &cobra.Command{
	Use:   "swagger",
	Short: `Gen swagger json docs by proto and api file`,
	RunE:  genswagger.Gen,
}

// genSdkCmd represents the gen sdk command
var genSdkCmd = &cobra.Command{
	Use:   "sdk",
	Short: `Generate sdk client by api file and proto file`,
	PreRun: func(_ *cobra.Command, _ []string) {
		if gensdk.Language == "ts" {
			console.Warning("[warning] ts client is still working...")
		}

		gensdk.Version = Version

		wd, err := os.Getwd()
		cobra.CheckErr(err)
		mod, err := mod.GetGoMod(wd)
		cobra.CheckErr(err)

		if gensdk.Output == "" {
			gensdk.Output = fmt.Sprintf("%s-%s", filepath.Base(mod.Path), gensdk.Language)
		}

		if gensdk.GoModule == "" {
			// module 为空, sdk 作为服务端的一个 package
			if gensdk.Language == "go" {
				gensdk.GoModule = filepath.ToSlash(filepath.Join(mod.Path, gensdk.Output))
			}
		} else {
			// module 不为空, 则生成 go.mod 文件
			gensdk.GenModule = true
		}

		if gensdk.GoPackage == "" {
			gensdk.GoPackage = strings.ReplaceAll(strings.ToLower(filepath.Base(gensdk.GoModule)), "-", "_")
		}

		if gensdk.Scope == "" {
			gensdk.Scope = filepath.Base(mod.Path)
			// 不支持 -, replace to _
			gensdk.Scope = strings.ReplaceAll(gensdk.Scope, "-", "_")
		}
	},
	RunE: gensdk.GenSdk,
}

// genDocsCmd represents the genDocs command
var genDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "jzero gen docs",
	Long:  `jzero gen docs`,
	PreRun: func(cmd *cobra.Command, args []string) {
		console.Warning("[warning] generate docs is still working...")
	},
	RunE: gendocs.Gen,
}

func init() {
	wd, _ := os.Getwd()

	{
		rootCmd.AddCommand(genCmd)

		// used for jzero
		genCmd.Flags().StringVarP(&embeded.Home, "home", "", filepath.Join(wd, ".template"), "set template home")
		genCmd.Flags().BoolVarP(&RemoveSuffix, "remove-suffix", "", true, "remove suffix Handler and Logic on filename or file content")
		genCmd.Flags().BoolVarP(&ChangeReplaceTypes, "change-replace-types", "", true, "if api file change, e.g. Request or Response type, change handler and logic file content types but not file")

		// used for goctl
		genCmd.Flags().StringVarP(&Style, "style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
		genCmd.Flags().StringSliceVarP(&ModelMysqlIgnoreColumns, "model-mysql-ignore-columns", "", []string{"create_at", "created_at", "create_time", "update_at", "updated_at", "update_time"}, "ignore columns of mysql model")
		genCmd.Flags().BoolVarP(&ModelMysqlDatasource, "model-mysql-datasource", "", false, "goctl model mysql datasource")
		genCmd.Flags().StringVarP(&ModelMysqlDatasourceUrl, "model-mysql-datasource-url", "", "", "goctl model mysql datasource url")
		genCmd.Flags().StringSliceVarP(&ModelMysqlDatasourceTable, "model-mysql-datasource-table", "", []string{"*"}, "goctl model mysql datasource table")
		genCmd.Flags().BoolVarP(&ModelMysqlCache, "model-mysql-cache", "", false, "goctl model mysql cache")
		genCmd.Flags().StringVarP(&ModelMysqlCachePrefix, "model-mysql-cache-prefix", "", "", "goctl model mysql cache prefix")
	}

	{
		genCmd.AddCommand(genSdkCmd)

		genSdkCmd.Flags().StringVarP(&gensdk.Language, "language", "l", "go", "set language")
		genSdkCmd.Flags().StringVarP(&gensdk.Output, "output", "o", "", "set output dir")
		genSdkCmd.Flags().StringVarP(&gensdk.GoModule, "goModule", "", "", "set module name")
		genSdkCmd.Flags().StringVarP(&gensdk.GoPackage, "goPackage", "", "", "set package name")
		genSdkCmd.Flags().StringVarP(&gensdk.ApiDir, "api-dir", "", filepath.Join("desc", "api"), "set input api dir")
		genSdkCmd.Flags().StringVarP(&gensdk.ProtoDir, "proto-dir", "", filepath.Join("desc", "proto"), "set input proto dir")
		genSdkCmd.Flags().BoolVarP(&gensdk.WrapResponse, "wrap-response", "", true, "warp response: code, data, message")
		genSdkCmd.Flags().StringVarP(&gensdk.Scope, "scope", "", "", "set scope name")
		genSdkCmd.Flags().StringVarP(&embeded.Home, "home", "", filepath.Join(wd, ".template"), "set template home")
	}

	{
		genCmd.AddCommand(genSwaggerCmd)

		genSwaggerCmd.Flags().StringVarP(&genswagger.Dir, "dir", "d", filepath.Join("desc", "swagger"), "set swagger output dir")
		genSwaggerCmd.Flags().StringVarP(&genswagger.ApiDir, "api-dir", "", filepath.Join("desc", "api"), "set input api dir")
		genSwaggerCmd.Flags().StringVarP(&genswagger.ProtoDir, "proto-dir", "", filepath.Join("desc", "proto"), "set input proto dir")
	}

	{
		genCmd.AddCommand(genDocsCmd)
	}

	{
		genCmd.AddCommand(genZRpcClientCmd)

		genZRpcClientCmd.Flags().StringVarP(&genzrpcclient.Style, "style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
		genZRpcClientCmd.Flags().StringVarP(&genzrpcclient.Output, "output", "o", "zrpcclient-go", "generate rpcclient code")
		genZRpcClientCmd.Flags().StringVarP(&genzrpcclient.Scope, "scope", "", "", "set scope name")
		genZRpcClientCmd.Flags().StringVarP(&genzrpcclient.GoModule, "goModule", "", "", "set go module name")
		genZRpcClientCmd.Flags().StringVarP(&genzrpcclient.GoPackage, "goPackage", "", "", "set package name")
	}
}
