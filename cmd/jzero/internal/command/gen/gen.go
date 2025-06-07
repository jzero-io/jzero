/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gen"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gendocs"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genswagger"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genzrpcclient"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/version"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: `Used to generate server code with api, proto, sql desc file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return gen.Run()
	},
	SilenceUsage: true,
}

// genZRpcClientCmd represents the rpcClient command
var genZRpcClientCmd = &cobra.Command{
	Use:   "zrpcclient",
	Short: `Gen zrpc client code by proto`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		cobra.CheckErr(err)
		mod, err := mod.GetGoMod(wd)
		cobra.CheckErr(err)
		if config.C.Gen.Zrpcclient.Scope == "" {
			config.C.Gen.Zrpcclient.Scope = filepath.Base(mod.Path)
			config.C.Gen.Zrpcclient.Scope = strings.ReplaceAll(config.C.Gen.Zrpcclient.Scope, "-", "_")
		}

		var genModule bool
		if config.C.Gen.Zrpcclient.GoModule == "" {
			config.C.Gen.Zrpcclient.GoModule = filepath.ToSlash(filepath.Join(mod.Path, config.C.Gen.Zrpcclient.Output))
		} else {
			genModule = true
		}

		if config.C.Gen.Zrpcclient.GoPackage == "" {
			config.C.Gen.Zrpcclient.GoPackage = strings.ReplaceAll(strings.ToLower(filepath.Base(config.C.Gen.Zrpcclient.GoModule)), "-", "_")
		}
		return genzrpcclient.Generate(genModule)
	},
}

// genSwaggerCmd represents the genSwagger command
var genSwaggerCmd = &cobra.Command{
	Use:   "swagger",
	Short: `Gen swagger json docs by proto and api file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genswagger.Gen()
	},
}

// genSdkCmd represents the gen sdk command
var genSdkCmd = &cobra.Command{
	Use:   "sdk",
	Short: `Generate sdk client by api file and proto file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.C.Gen.Sdk.Language == "ts" {
			console.Warning("[warning] ts client is still working...")
		}

		mod, err := mod.GetGoMod(config.C.Wd())
		if err != nil {
			return err
		}

		if config.C.Gen.Sdk.Output == "" {
			config.C.Gen.Sdk.Output = fmt.Sprintf("%s-%s", filepath.Base(mod.Path), config.C.Gen.Sdk.Language)
		}

		var genModule bool
		if config.C.Gen.Sdk.GoModule == "" {
			// module 为空, sdk 作为服务端的一个 package
			if config.C.Gen.Sdk.Language == "go" {
				config.C.Gen.Sdk.GoModule = filepath.ToSlash(filepath.Join(mod.Path, config.C.Gen.Sdk.Output))
			}
		} else {
			// module 不为空, 则生成 go.mod 文件
			if !config.C.Gen.Sdk.Mono {
				genModule = true
			}
		}

		if config.C.Gen.Sdk.GoPackage == "" {
			config.C.Gen.Sdk.GoPackage = strings.ReplaceAll(strings.ToLower(filepath.Base(config.C.Gen.Sdk.GoModule)), "-", "_")
		}

		if config.C.Gen.Sdk.Scope == "" {
			config.C.Gen.Sdk.Scope = filepath.Base(mod.Path)
			// 不支持 -, replace to _
			config.C.Gen.Sdk.Scope = strings.ReplaceAll(config.C.Gen.Sdk.Scope, "-", "_")
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		if embeded.Home == "" {
			embeded.Home = filepath.Join(homeDir, ".jzero", "templates", version.Version)
		}
		return gensdk.GenSdk(genModule)
	},
}

// genDocsCmd represents the genDocs command
var genDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "jzero gen docs",
	Long:  `jzero gen docs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config.C.Gen.Docs.Output = filepath.Join("desc", "docs", config.C.Gen.Docs.Format)
		return gendocs.Gen()
	},
	Aliases: []string{"doc"},
}

func GetCommand() *cobra.Command {
	{
		genCmd.PersistentFlags().StringP("style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
		genCmd.PersistentFlags().StringP("home", "", ".template", "set template home")
		genCmd.PersistentFlags().StringSliceP("desc", "", []string{}, "set desc path")
		genCmd.PersistentFlags().StringSliceP("desc-ignore", "", []string{}, "set desc ignore path")

		genCmd.Flags().BoolP("git-change", "", false, "set is git change, if changes then generate code")
		genCmd.Flags().BoolP("route2code", "", false, "is generate route2code")
		genCmd.Flags().BoolP("rpc-client", "", false, "is generate rpc client code by goctl")

		// common model, support more db
		genCmd.Flags().StringP("model-driver", "", "mysql", "goctl model driver. mysql or postgres")
		genCmd.Flags().BoolP("model-strict", "", false, "goctl model strict mode, see [https://go-zero.dev/docs/tutorials/cli/model]")
		genCmd.Flags().StringSliceP("model-ignore-columns", "", []string{"create_at", "created_at", "create_time", "update_at", "updated_at", "update_time"}, "ignore columns of mysql model")
		genCmd.Flags().StringP("model-schema", "", "", "model schema")
		genCmd.Flags().BoolP("model-datasource", "", false, "goctl datasource")
		genCmd.Flags().StringSliceP("model-datasource-url", "", []string{}, "goctl model datasource url")
		genCmd.Flags().StringSliceP("model-datasource-table", "", []string{"*"}, "goctl model datasource table")
		genCmd.Flags().BoolP("model-cache", "", false, "goctl model cache")
		genCmd.Flags().StringP("model-cache-prefix", "", "cache", "goctl model cache prefix")
		genCmd.Flags().BoolP("model-create-table-ddl", "", false, "is generate create table ddl, only datasource mode takes effective")
	}

	{
		genCmd.AddCommand(genSdkCmd)

		genSdkCmd.Flags().StringP("language", "l", "go", "set language")
		genSdkCmd.Flags().StringP("output", "o", "", "set output dir")
		genSdkCmd.Flags().StringP("goModule", "", "", "set go module name")
		genSdkCmd.Flags().StringP("goVersion", "", "", "set go version, only effect when having goModule flag")
		genSdkCmd.Flags().StringP("goPackage", "", "", "set package name")
		genSdkCmd.Flags().BoolP("wrap-response", "", true, "wrap response: code, data, message")
		genSdkCmd.Flags().StringP("scope", "", "", "set scope name")
		genSdkCmd.Flags().BoolP("mono", "", false, "mono sdk project under go mod project")
	}

	{
		genCmd.AddCommand(genSwaggerCmd)

		genSwaggerCmd.Flags().StringP("output", "o", filepath.Join("desc", "swagger"), "set swagger output dir")
		genSwaggerCmd.Flags().BoolP("route2code", "", false, "is generate route2code")
		genSwaggerCmd.Flags().BoolP("merge", "", false, "is merge muti swagger to one file, goctl version >= v1.8.3 available")
	}

	{
		genCmd.AddCommand(genDocsCmd)

		genDocsCmd.Flags().StringP("output", "o", filepath.Join("desc", "docs", "md"), "set docs output dir")
		genDocsCmd.Flags().StringP("format", "", "md", "set output format")
	}

	{
		genCmd.AddCommand(genZRpcClientCmd)

		genZRpcClientCmd.Flags().StringP("output", "o", "zrpcclient-go", "generate rpcclient code")
		genZRpcClientCmd.Flags().StringP("pb-dir", "", "", "set output pb dir ")
		genZRpcClientCmd.Flags().StringP("client-dir", "", "", "set output client dir ")
		genZRpcClientCmd.Flags().StringP("scope", "", "", "set scope name")
		genZRpcClientCmd.Flags().StringP("goModule", "", "", "set go module name")
		genZRpcClientCmd.Flags().StringP("goVersion", "", "", "set go version, only effect when having goModule flag")
		genZRpcClientCmd.Flags().StringP("goPackage", "", "", "set package name")
	}

	return genCmd
}
