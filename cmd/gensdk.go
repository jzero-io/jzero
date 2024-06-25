/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/pkg/mod"

	"github.com/jzero-io/jzero/internal/gen/gensdk"

	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
)

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

func init() {
	genSdkCmd.Flags().StringVarP(&gensdk.Language, "language", "l", "go", "set language")
	genSdkCmd.Flags().StringVarP(&gensdk.Dir, "dir", "d", "", "set dir")

	genSdkCmd.Flags().StringVarP(&gensdk.WorkingDir, "working-dir", "w", "", "set working dir")

	genSdkCmd.Flags().StringVarP(&gensdk.Module, "module", "m", "", "set module name")

	genSdkCmd.Flags().StringVarP(&gensdk.ApiDir, "api-dir", "", filepath.Join("desc", "api"), "set input api dir")
	genSdkCmd.Flags().StringVarP(&gensdk.ProtoDir, "proto-dir", "", filepath.Join("desc", "proto"), "set input proto dir")
	genSdkCmd.Flags().BoolVarP(&gensdk.WarpResponse, "warp-response", "", false, "warp response: code, data, message")
	genSdkCmd.Flags().StringVarP(&gensdk.Scope, "scope", "", "", "set scope name")
	genSdkCmd.Flags().StringVarP(&embeded.Home, "home", "", "", "set template home")
}
