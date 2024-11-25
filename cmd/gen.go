/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/internal/gen/gendocs"
	"github.com/jzero-io/jzero/internal/gen/gensdk"
	"github.com/jzero-io/jzero/internal/gen/genswagger"
	"github.com/jzero-io/jzero/internal/gen/genzrpcclient"
	"github.com/jzero-io/jzero/pkg"
	"github.com/jzero-io/jzero/pkg/mod"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: `Used to generate server/client code`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(config.C.Gen.Hooks.Before) > 0 {
			fmt.Printf("%s hooks \n", color.WithColor("Start Before", color.FgGreen))
		}
		for _, v := range config.C.Gen.Hooks.Before {
			fmt.Printf("%s command %s\n", color.WithColor("Run", color.FgGreen), v)
			err := pkg.Run(v, config.C.Wd())
			if err != nil {
				return err
			}
		}
		if len(config.C.Gen.Hooks.Before) > 0 {
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := os.UserHomeDir()
		config.C.Gen.Home, _ = homedir.Expand(config.C.Gen.Home)
		if !pathx.FileExists(config.C.Gen.Home) {
			config.C.Gen.Home = filepath.Join(home, ".jzero", Version)
		}
		embeded.Home = config.C.Gen.Home
		return gen.Run(config.C)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		if len(config.C.Gen.Hooks.After) > 0 {
			fmt.Printf("%s hooks \n", color.WithColor("Start After", color.FgGreen))
		}
		for _, v := range config.C.Gen.Hooks.After {
			fmt.Printf("%s command %s\n", color.WithColor("Run", color.FgGreen), v)
			err := pkg.Run(v, config.C.Wd())
			if err != nil {
				return err
			}
		}
		if len(config.C.Gen.Hooks.After) > 0 {
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		}
		return nil
	},
	SilenceUsage: true,
}

// genZRpcClientCmd represents the rpcClient command
var genZRpcClientCmd = &cobra.Command{
	Use:   "zrpcclient",
	Short: `Gen zrpc client code by proto`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(config.C.Gen.Zrpcclient.Hooks.Before) > 0 {
			fmt.Printf("%s hooks \n", color.WithColor("Start Before", color.FgGreen))
		}
		for _, v := range config.C.Gen.Zrpcclient.Hooks.Before {
			fmt.Printf("%s command %s\n", color.WithColor("Run", color.FgGreen), v)
			err := pkg.Run(v, config.C.Wd())
			if err != nil {
				return err
			}
		}
		if len(config.C.Gen.Zrpcclient.Hooks.Before) > 0 {
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		}
		return nil
	},
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
		return genzrpcclient.Generate(config.C, genModule)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		if len(config.C.Gen.Zrpcclient.Hooks.After) > 0 {
			fmt.Printf("%s hooks \n", color.WithColor("Start After", color.FgGreen))
		}
		for _, v := range config.C.Gen.Zrpcclient.Hooks.After {
			fmt.Printf("%s command %s\n", color.WithColor("Run", color.FgGreen), v)
			err := pkg.Run(v, config.C.Wd())
			if err != nil {
				return err
			}
		}
		if len(config.C.Gen.Zrpcclient.Hooks.After) > 0 {
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		}
		return nil
	},
}

// genSwaggerCmd represents the genSwagger command
var genSwaggerCmd = &cobra.Command{
	Use:   "swagger",
	Short: `Gen swagger json docs by proto and api file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genswagger.Gen(config.C)
	},
}

// genSdkCmd represents the gen sdk command
var genSdkCmd = &cobra.Command{
	Use:   "sdk",
	Short: `Generate sdk client by api file and proto file`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(config.C.Gen.Sdk.Hooks.Before) > 0 {
			fmt.Printf("%s hooks \n", color.WithColor("Start Before", color.FgGreen))
		}
		for _, v := range config.C.Gen.Sdk.Hooks.Before {
			fmt.Printf("%s command %s\n", color.WithColor("Run", color.FgGreen), v)
			err := pkg.Run(v, config.C.Wd())
			if err != nil {
				return err
			}
		}
		if len(config.C.Gen.Sdk.Hooks.Before) > 0 {
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		}
		return nil
	},
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
			genModule = true
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
			embeded.Home = filepath.Join(homeDir, ".jzero", Version)
		}
		return gensdk.GenSdk(config.C, genModule)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		if len(config.C.Gen.Sdk.Hooks.After) > 0 {
			fmt.Printf("%s hooks \n", color.WithColor("Start After", color.FgGreen))
		}
		for _, v := range config.C.Gen.Sdk.Hooks.After {
			fmt.Printf("%s command %s\n", color.WithColor("Run", color.FgGreen), v)
			err := pkg.Run(v, config.C.Wd())
			if err != nil {
				return err
			}
		}
		if len(config.C.Gen.Sdk.Hooks.After) > 0 {
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		}
		return nil
	},
}

// genDocsCmd represents the genDocs command
var genDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "jzero gen docs",
	Long:  `jzero gen docs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config.C.Gen.Docs.Output = filepath.Join("desc", "docs", config.C.Gen.Docs.Format)
		return gendocs.Gen(config.C)
	},
	Aliases: []string{"doc"},
}

func init() {
	wd, _ := os.Getwd()

	{
		rootCmd.AddCommand(genCmd)

		// used for jzero
		genCmd.Flags().BoolP("change-logic-types", "", false, "if api file or proto change, e.g. Request or Response type, change logic file content types but not file")
		genCmd.Flags().BoolP("regen-api-handler", "", false, "")

		// fix rpc style
		genCmd.Flags().BoolP("rpc-style-patch", "", false, "")
		// git-diff
		genCmd.Flags().BoolP("git-change", "", false, "set is git change, if changes then generate code")
		genCmd.Flags().StringP("api-git-change-path", "", filepath.Join("desc", "api"), "set api git change path")
		genCmd.Flags().StringP("model-git-change-path", "", filepath.Join("desc", "sql"), "set model git change path")
		genCmd.Flags().StringP("proto-git-change-path", "", filepath.Join("desc", "proto"), "set proto git change path")

		genCmd.Flags().StringSliceP("desc", "", []string{}, "set desc path")
		genCmd.Flags().StringSliceP("desc-ignore", "", []string{}, "set desc ignore path")

		genCmd.Flags().BoolP("rpc-client", "", false, "is generate rpc client code by goctl")

		// used for goctl
		// gen command persistentFlags
		genCmd.PersistentFlags().StringP("style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
		genCmd.PersistentFlags().StringP("home", "", filepath.Join(wd, ".template"), "set template home")
		genCmd.PersistentFlags().BoolP("split-api-types-dir", "", false, "")

		genCmd.Flags().BoolP("model-mysql-strict", "", false, "goctl model mysql strict mode, see [https://go-zero.dev/docs/tutorials/cli/model]")
		genCmd.Flags().StringSliceP("model-mysql-ignore-columns", "", []string{"create_at", "created_at", "create_time", "update_at", "updated_at", "update_time"}, "ignore columns of mysql model")
		genCmd.Flags().StringP("model-mysql-ddl-database", "", "", "goctl model mysql ddl database")
		genCmd.Flags().BoolP("model-mysql-datasource", "", false, "goctl model mysql datasource")
		genCmd.Flags().StringP("model-mysql-datasource-url", "", "", "goctl model mysql datasource url")
		genCmd.Flags().StringSliceP("model-mysql-datasource-table", "", []string{"*"}, "goctl model mysql datasource table")
		genCmd.Flags().BoolP("model-mysql-cache", "", false, "goctl model mysql cache")
		genCmd.Flags().StringP("model-mysql-cache-prefix", "", "", "goctl model mysql cache prefix")
		genCmd.Flags().BoolP("gen-mysql-create-table-ddl", "", false, "is generate mysql create table ddl, only datasource mode takes effective")
	}

	{
		genCmd.AddCommand(genSdkCmd)

		genSdkCmd.Flags().StringP("language", "l", "go", "set language")
		genSdkCmd.Flags().StringP("output", "o", "", "set output dir")
		genSdkCmd.Flags().StringP("goModule", "", "", "set go module name")
		genSdkCmd.Flags().StringP("goVersion", "", "", "set go version, only effect when having goModule flag")
		genSdkCmd.Flags().StringP("goPackage", "", "", "set package name")
		genSdkCmd.Flags().StringP("api-dir", "", filepath.Join("desc", "api"), "set input api dir")
		genSdkCmd.Flags().StringP("proto-dir", "", filepath.Join("desc", "proto"), "set input proto dir")
		genSdkCmd.Flags().BoolP("wrap-response", "", true, "wrap response: code, data, message")
		genSdkCmd.Flags().StringP("scope", "", "", "set scope name")
	}

	{
		genCmd.AddCommand(genSwaggerCmd)

		genSwaggerCmd.Flags().StringP("output", "o", filepath.Join("desc", "swagger"), "set swagger output dir")
		genSwaggerCmd.Flags().StringP("api-dir", "", filepath.Join("desc", "api"), "set input api dir")
		genSwaggerCmd.Flags().StringP("proto-dir", "", filepath.Join("desc", "proto"), "set input proto dir")
	}

	{
		genCmd.AddCommand(genDocsCmd)

		genDocsCmd.Flags().StringP("output", "o", filepath.Join("desc", "docs", "md"), "set docs output dir")
		genDocsCmd.Flags().StringP("format", "", "md", "set output format")
		genDocsCmd.Flags().StringP("api-dir", "", filepath.Join("desc", "api"), "set input api dir")
		genDocsCmd.Flags().StringP("proto-dir", "", filepath.Join("desc", "proto"), "set input proto dir")
	}

	{
		genCmd.AddCommand(genZRpcClientCmd)

		genZRpcClientCmd.Flags().StringP("output", "o", "zrpcclient-go", "generate rpcclient code")
		genZRpcClientCmd.Flags().StringP("scope", "", "", "set scope name")
		genZRpcClientCmd.Flags().StringP("goModule", "", "", "set go module name")
		genZRpcClientCmd.Flags().StringP("goVersion", "", "", "set go version, only effect when having goModule flag")
		genZRpcClientCmd.Flags().StringP("goPackage", "", "", "set package name")
	}
}
