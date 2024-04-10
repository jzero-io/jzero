/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"bytes"
	"fmt"
	"github.com/jaronnie/jzero/config"
	"github.com/jaronnie/jzero/protosets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "jzerod init",
	Long:  `jzerod init`,
	Run: func(cmd *cobra.Command, args []string) {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Check $HOME/.jzero/protosets
		err = os.MkdirAll(filepath.Join(home, ".jzero", "protosets"), 0755)
		cobra.CheckErr(err)

		file, err := config.Config.ReadFile("config.toml")
		cobra.CheckErr(err)

		// write protosets
		dir, err := protosets.Protosets.ReadDir(".")
		cobra.CheckErr(err)
		for _, d := range dir {
			pb, err := protosets.Protosets.ReadFile(d.Name())
			cobra.CheckErr(err)
			err = os.WriteFile(filepath.Join(home, ".jzero", "protosets", d.Name()), pb, 0644)
			cobra.CheckErr(err)
		}

		v := viper.New()
		v.SetConfigType("toml")
		err = v.ReadConfig(bytes.NewBuffer(file))
		cobra.CheckErr(err)

		oldProtoSets := v.GetStringSlice("Gateway.Upstreams.0.ProtoSets")
		var newProtoSets []string

		for _, protoSet := range oldProtoSets {
			newProtoSets = append(newProtoSets, filepath.Join(home, ".jzero", protoSet))
		}
		v.Set("Gateway.Upstreams.0.ProtoSets", newProtoSets)

		err = v.WriteConfigAs(filepath.Join(home, ".jzero", "config.toml"))
		cobra.CheckErr(err)

		fmt.Println("init success")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
