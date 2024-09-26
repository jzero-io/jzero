/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: `Check and install all needed tools`,
	Run: func(cmd *cobra.Command, args []string) {
		log := console.NewColorConsole(true)

		log.Info("[jzero-env]: looking up goctl")
		// install goctl
		_, err := env.LookPath("goctl")
		if err != nil {
			log.Warning(`[jzero-env]: goctl is not found in PATH`)
			err = golang.Install("github.com/zeromicro/go-zero/tools/goctl@latest")
			cobra.CheckErr(err)
		}
		if _, err = env.LookPath("goctl"); err == nil {
			log.Success(`[jzero-env]: "goctl" is installed`)
		} else {
			log.Fatalln("[jzero-env]: env check failed, goctl is not installed")
		}

		// goctl env check
		resp, err := execx.Run("goctl env check --install --verbose --force", "")
		cobra.CheckErr(err)
		fmt.Println(resp)

		log.Info("\n[jzero-env]: looking up protoc-gen-go-grpc")
		_, err = env.LookPath("protoc-gen-go-grpc")
		if err != nil {
			_ = golang.Install("google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
		}
		if _, err = env.LookPath("protoc-gen-go-grpc"); err == nil {
			log.Success(`[jzero-env]: "protoc-gen-go-grpc" is installed`)
		} else {
			log.Fatalln("[jzero-env] env check failed, protoc-gen-go-grpc is not installed")
		}

		log.Info("\n[jzero-env]: looking up goctl-swagger")
		_, err = env.LookPath("goctl-swagger")
		if err != nil {
			_ = golang.Install("github.com/jzero-io/goctl-swagger@latest")
		}
		if _, err = env.LookPath("goctl-swagger"); err == nil {
			log.Success(`[jzero-env]: "goctl-swagger" is installed`)
		} else {
			log.Warning("[jzero-env] warning: env check failed, goctl-swagger is not installed")
		}

		log.Info("\n[jzero-env]: looking up protoc-gen-openapiv2")
		_, err = env.LookPath("protoc-gen-openapiv2")
		if err != nil {
			_ = golang.Install("github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest")
		}
		if _, err = env.LookPath("protoc-gen-openapiv2"); err == nil {
			log.Success(`[jzero-env]: "protoc-gen-openapiv2" is installed`)
		} else {
			log.Warning("[jzero-env] warning: env check failed, protoc-gen-openapiv2 is not installed")
		}

		// protoc validator
		log.Info("\n[jzero-env]: looking up protoc-gen-validate")
		_, err = env.LookPath("protoc-gen-validate")
		if err != nil {
			_ = golang.Install("github.com/jzero-io/protoc-gen-validate@latest")
		}
		if _, err = env.LookPath("protoc-gen-validate"); err == nil {
			log.Success(`[jzero-env]: "protoc-gen-validate" is installed`)
		} else {
			log.Warning("[jzero-env] warning: env check failed, protoc-gen-validate is not installed")
		}

		// protoc-gen-doc
		log.Info("\n[jzero-env]: looking up protoc-gen-doc")
		_, err = env.LookPath("protoc-gen-doc")
		if err != nil {
			_ = golang.Install("github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest")
		}
		if _, err = env.LookPath("protoc-gen-doc"); err == nil {
			log.Success(`[jzero-env]: "protoc-gen-doc" is installed`)
		} else {
			log.Warning("[jzero-env] warning: env check failed, protoc-gen-doc is not installed")
		}

		log.Success("\n[jzero-env]: congratulations! your jzero environment is ready!")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
