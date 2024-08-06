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
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.SetConfig(CfgFile, cmd.Parent().Use+"."+cmd.Use, cmd.Flags())
		if err != nil {
			return err
		}

		if c.Gen.Sdk.Language == "ts" {
			console.Warning("[warning] ts client is still working...")
		}

		mod, err := mod.GetGoMod(c.Gen.Wd())
		cobra.CheckErr(err)

		if c.Gen.Sdk.Output == "" {
			c.Gen.Sdk.Output = fmt.Sprintf("%s-%s", filepath.Base(mod.Path), c.Gen.Sdk.Language)
		}

		var genModule bool
		if c.Gen.Sdk.GoModule == "" {
			// module 为空, sdk 作为服务端的一个 package
			if c.Gen.Sdk.Language == "go" {
				c.Gen.Sdk.GoModule = filepath.ToSlash(filepath.Join(mod.Path, c.Gen.Sdk.Output))
			}
		} else {
			// module 不为空, 则生成 go.mod 文件
			genModule = true
		}

		if c.Gen.Sdk.GoPackage == "" {
			c.Gen.Sdk.GoPackage = strings.ReplaceAll(strings.ToLower(filepath.Base(c.Gen.Sdk.GoModule)), "-", "_")
		}

		if c.Gen.Sdk.Scope == "" {
			c.Gen.Sdk.Scope = filepath.Base(mod.Path)
			// 不支持 -, replace to _
			c.Gen.Sdk.Scope = strings.ReplaceAll(c.Gen.Sdk.Scope, "-", "_")
		}

		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		if embeded.Home == "" {
			embeded.Home = filepath.Join(homeDir, ".jzero", Version)
		}
		return gensdk.GenSdk(c.Gen.Sdk, genModule)
	},
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
		genCmd.Flags().StringP("home", "", filepath.Join(wd, ".template"), "set template home")
		genCmd.Flags().BoolP("remove-suffix", "", true, "remove suffix Handler and Logic on filename or file content")
		genCmd.Flags().BoolP("change-replace-types", "", true, "if api file change, e.g. Request or Response type, change handler and logic file content types but not file")

		// used for goctl
		genCmd.Flags().StringP("style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
		genCmd.Flags().StringSliceP("model-mysql-ignore-columns", "", []string{"create_at", "created_at", "create_time", "update_at", "updated_at", "update_time"}, "ignore columns of mysql model")
		genCmd.Flags().BoolP("model-mysql-datasource", "", false, "goctl model mysql datasource")
		genCmd.Flags().StringP("model-mysql-datasource-url", "", "", "goctl model mysql datasource url")
		genCmd.Flags().StringSliceP("model-mysql-datasource-table", "", []string{"*"}, "goctl model mysql datasource table")
		genCmd.Flags().BoolP("model-mysql-cache", "", false, "goctl model mysql cache")
		genCmd.Flags().StringP("model-mysql-cache-prefix", "", "", "goctl model mysql cache prefix")
	}

	{
		genCmd.AddCommand(genSdkCmd)

		genSdkCmd.Flags().StringP("language", "l", "go", "set language")
		genSdkCmd.Flags().StringP("output", "o", "", "set output dir")
		genSdkCmd.Flags().StringP("goModule", "", "", "set module name")
		genSdkCmd.Flags().StringP("goPackage", "", "", "set package name")
		genSdkCmd.Flags().StringP("api-dir", "", filepath.Join("desc", "api"), "set input api dir")
		genSdkCmd.Flags().StringP("proto-dir", "", filepath.Join("desc", "proto"), "set input proto dir")
		genSdkCmd.Flags().BoolP("wrap-response", "", true, "warp response: code, data, message")
		genSdkCmd.Flags().StringP("scope", "", "", "set scope name")
		genSdkCmd.Flags().StringP("home", "", filepath.Join(wd, ".template"), "set template home")
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
