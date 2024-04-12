/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "jzero gen code",
	Long:  `jzero gen code`,
	RunE:  gen,
}

func gen(cmd *cobra.Command, args []string) error {
	time.Sleep(time.Second * 15)
	wd, err := os.Getwd()
	cobra.CheckErr(err)
	// read proto dir
	ds, err := os.ReadDir(filepath.Join(wd, "daemon", "proto"))
	cobra.CheckErr(err)
	for _, v := range ds {
		if v.IsDir() {
			continue
		}
		if strings.HasSuffix(v.Name(), "proto") {
			command := fmt.Sprintf("goctl rpc protoc daemon/proto/%s  -I./daemon/proto --go_out=./daemon --go-grpc_out=./daemon  --zrpc_out=./daemon --client=false -m", v.Name())
			_, err = Run(command, wd)
			cobra.CheckErr(err)

			fileBase := v.Name()[0 : len(v.Name())-len(path.Ext(v.Name()))]
			_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", fileBase)))
		}
	}
	_ = os.RemoveAll(filepath.Join(wd, "daemon", "etc"))

	// read api file
	v := viper.New()
	v.SetConfigFile(filepath.Join(wd, "config.toml"))
	v.SetConfigType("toml")
	err = v.ReadInConfig()
	cobra.CheckErr(err)
	command := fmt.Sprintf("goctl api go --api daemon/api/%s.api --dir ./daemon", v.GetString("APP"))
	_, err = Run(command, wd)
	cobra.CheckErr(err)
	_ = os.Remove(filepath.Join(wd, "daemon", fmt.Sprintf("%s.go", v.Get("APP"))))
	return nil
}

func init() {
	rootCmd.AddCommand(genCmd)
}
