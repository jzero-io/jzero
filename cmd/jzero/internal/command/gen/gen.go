/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package gen

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rinchsan/gosimports"
	"github.com/spf13/cobra"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gen"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genswagger"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genzrpcclient"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
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

		var genModule bool
		if config.C.Gen.Zrpcclient.GoModule == "" {
			// module 为空, zrpcclient 作为服务端的一个 package
			output := config.C.Gen.Zrpcclient.Output

			// 计算输出路径的绝对路径
			absOutput, err := filepath.Abs(output)
			cobra.CheckErr(err)

			// 计算相对于 go.mod 所在目录的相对路径
			relPath, err := filepath.Rel(mod.Dir, absOutput)
			cobra.CheckErr(err)

			config.C.Gen.Zrpcclient.GoModule = filepath.ToSlash(filepath.Join(mod.Path, relPath))
		} else {
			genModule = true
		}
		gosimports.LocalPrefix = config.C.Gen.Zrpcclient.GoModule

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

func GetCommand() *cobra.Command {
	{
		genCmd.Flags().StringSliceP("desc", "", []string{}, "set desc path")
		genCmd.Flags().StringSliceP("desc-ignore", "", []string{}, "set desc ignore path")
		genCmd.Flags().BoolP("git-change", "", false, "set is git change, if changes then generate code")
		genCmd.Flags().BoolP("route2code", "", false, "is generate route2code")
		genCmd.Flags().StringSliceP("proto-include", "", []string{}, "proto include path")
		genCmd.Flags().StringP("model-driver", "", "mysql", "goctl model driver. mysql or postgres")
		genCmd.Flags().BoolP("model-strict", "", false, "goctl model strict mode, see [https://go-zero.dev/docs/tutorials/cli/model]")
		genCmd.Flags().StringSliceP("model-ignore-columns", "", []string{"create_at", "created_at", "create_time", "update_at", "updated_at", "update_time"}, "ignore columns of mysql model")
		genCmd.Flags().StringP("model-schema", "", "", "model schema")
		genCmd.Flags().BoolP("model-datasource", "", false, "goctl datasource")
		genCmd.Flags().StringSliceP("model-datasource-url", "", []string{}, "goctl model datasource url")
		genCmd.Flags().StringSliceP("model-datasource-table", "", []string{"*"}, "goctl model datasource table")
		genCmd.Flags().BoolP("model-cache", "", false, "goctl model cache")
		genCmd.Flags().StringSliceP("model-cache-table", "", []string{"*"}, "goctl model cache tables")
		genCmd.Flags().StringP("model-cache-prefix", "", "cache", "goctl model cache prefix")
		genCmd.Flags().StringSliceP("mongo-type", "", []string{}, "mongo type name")
		genCmd.Flags().BoolP("mongo-cache", "", false, " Generate code with cache prefix [optional]")
		genCmd.Flags().StringP("mongo-cache-prefix", "", "cache", "mongo cache prefix")
		genCmd.Flags().StringSliceP("mongo-cache-type", "", []string{"*"}, "mongo cache type names to enable caching")
		genCmd.Flags().BoolP("rpc-client", "", false, "generate rpc client code")
	}

	{
		genCmd.AddCommand(genSwaggerCmd)

		genSwaggerCmd.Flags().StringSliceP("desc", "", []string{}, "set desc path")
		genSwaggerCmd.Flags().StringSliceP("desc-ignore", "", []string{}, "set desc ignore path")
		genSwaggerCmd.Flags().StringP("output", "o", filepath.Join("desc", "swagger"), "set swagger output dir")
		genSwaggerCmd.Flags().BoolP("route2code", "", false, "is generate route2code")
		genSwaggerCmd.Flags().BoolP("merge", "", true, "is merge muti swagger to one file")
	}

	{
		genCmd.AddCommand(genZRpcClientCmd)

		genZRpcClientCmd.Flags().StringSliceP("desc", "", []string{}, "set desc path")
		genZRpcClientCmd.Flags().StringSliceP("desc-ignore", "", []string{}, "set desc ignore path")
		genZRpcClientCmd.Flags().StringSliceP("proto-include", "", []string{}, "proto include path")
		genZRpcClientCmd.Flags().StringP("output", "o", "zrpcclient-go", "generate rpcclient code")
		genZRpcClientCmd.Flags().StringP("goModule", "", "", "set go module name")
		genZRpcClientCmd.Flags().StringP("goVersion", "", "", "set go version, only effect when having goModule flag")
		genZRpcClientCmd.Flags().StringP("goPackage", "", "", "set package name")
	}

	return genCmd
}
