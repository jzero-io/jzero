/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/internal/gen/gendocs"
	"github.com/jzero-io/jzero/internal/gen/gensdk"
	"github.com/jzero-io/jzero/pkg/mod"

	"github.com/jzero-io/jzero/internal/gen/genswagger"

	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "jzero gen code",
	Long:  `jzero gen code`,
	PreRun: func(_ *cobra.Command, _ []string) {
		gen.Version = Version
		gen.AppDir = strings.TrimPrefix(gen.AppDir, ".")

		// check go-zero api template
		home, _ := os.UserHomeDir()
		if !pathx.FileExists(filepath.Join(home, ".jzero", Version, "go-zero")) {
			err := embeded.WriteTemplateDir(filepath.Join("go-zero"), filepath.Join(home, ".jzero", Version, "go-zero"))
			cobra.CheckErr(err)
		}
	},
	RunE:         gen.Gen,
	SilenceUsage: false,
}

// genSwaggerCmd represents the genSwagger command
var genSwaggerCmd = &cobra.Command{
	Use:   "swagger",
	Short: "jzero gen swagger",
	Long:  `jzero gen swagger`,
	RunE:  genswagger.Gen,
}

// genSdkCmd represents the gen sdk command
var genSdkCmd = &cobra.Command{
	Use:   "sdk",
	Short: "jzero gensdk",
	Long:  `jzero gensdk. Generate sdk client by api file and proto file`,
	PreRun: func(_ *cobra.Command, _ []string) {
		gensdk.Version = Version

		wd, err := os.Getwd()
		cobra.CheckErr(err)
		mod, err := mod.GetGoMod(wd)
		cobra.CheckErr(err)

		if gensdk.Module == "" {
			gensdk.Module = fmt.Sprintf("%s-go", mod.Path)
		}

		if gensdk.Dir == "" {
			gensdk.Dir = fmt.Sprintf("%s-go", mod.Path)
		}

		if gensdk.Scope == "" {
			gensdk.Scope = filepath.Base(mod.Path)
		}
	},
	RunE: gensdk.GenSdk,
}

// genDocsCmd represents the genDocs command
var genDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "jzero gen docs",
	Long:  `jzero gen docs`,
	RunE:  gendocs.Gen,
}

func init() {
	{
		rootCmd.AddCommand(genCmd)

		genCmd.Flags().StringVarP(&gen.WorkingDir, "working-dir", "w", "", "set working dir")
		genCmd.Flags().StringVarP(&gen.AppDir, "app-dir", "", ".", "set app dir")
		dir, _ := os.UserHomeDir()
		genCmd.Flags().StringVarP(&embeded.Home, "home", "", filepath.Join(dir, ".jzero"), "set template home")
		genCmd.Flags().StringVarP(&gen.Style, "style", "", "gozero", "The file naming format, see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
		genCmd.Flags().BoolVarP(&gen.RemoveSuffix, "remove-suffix", "", false, "remove suffix Handler and Logic on filename or file content")
		genCmd.Flags().BoolVarP(&gen.ChangeReplaceTypes, "change-replace-types", "", false, "if api file change, e.g. Request or Response type, change handler and logic file content types but not file")

		genCmd.Flags().StringSliceVarP(&gen.ModelMysqlIgnoreColumns, "model-mysql-ignore-columns", "", []string{"create_at", "created_at", "create_time", "update_at", "updated_at", "update_time"}, "ignore columns of mysql model")
		genCmd.Flags().BoolVarP(&gen.ModelMysqlDatasource, "model-mysql-datasource", "", false, "goctl model mysql datasource")
		genCmd.Flags().StringVarP(&gen.ModelMysqlDatasourceUrl, "model-mysql-datasource-url", "", "", "goctl model mysql datasource url")
		genCmd.Flags().StringSliceVarP(&gen.ModelMysqlDatasourceTable, "model-mysql-datasource-table", "", []string{"*"}, "goctl model mysql datasource table")
		genCmd.Flags().BoolVarP(&gen.ModelMysqlCache, "model-mysql-cache", "", false, "goctl model mysql cache")
		genCmd.Flags().StringVarP(&gen.ModelMysqlCachePrefix, "model-mysql-cache-prefix", "", "", "goctl model mysql cache prefix")
	}

	{
		genCmd.AddCommand(genSdkCmd)

		genSdkCmd.Flags().StringVarP(&gensdk.Language, "language", "l", "go", "set language")
		genSdkCmd.Flags().StringVarP(&gensdk.Dir, "dir", "d", "", "set dir")
		genSdkCmd.Flags().StringVarP(&gensdk.WorkingDir, "working-dir", "w", "", "set working dir")
		genSdkCmd.Flags().StringVarP(&gensdk.Module, "module", "m", "", "set module name")
		genSdkCmd.Flags().StringVarP(&gensdk.ApiDir, "api-dir", "", filepath.Join("desc", "api"), "set input api dir")
		genSdkCmd.Flags().StringVarP(&gensdk.ProtoDir, "proto-dir", "", filepath.Join("desc", "proto"), "set input proto dir")
		genSdkCmd.Flags().BoolVarP(&gensdk.WrapResponse, "wrap-response", "", false, "warp response: code, data, message")
		genSdkCmd.Flags().StringVarP(&gensdk.Scope, "scope", "", "", "set scope name")
		genSdkCmd.Flags().StringVarP(&embeded.Home, "home", "", "", "set template home")
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
}
