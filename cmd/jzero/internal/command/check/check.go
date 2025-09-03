/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package check

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"

	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
)

// RunCheck executes the check logic and can be called from other places
func RunCheck(verbose bool) error {
	log := console.NewColorConsole(true)

	// Detect frame type
	frameType := desc.GetFrameType()
	if frameType == "" {
		return nil
	}

	// install goctl
	_, err := env.LookPath("goctl")
	if err != nil {
		fmt.Printf("%s goctl\n", color.WithColor("Installing tool", color.FgGreen))
		err = golang.Install("github.com/zeromicro/go-zero/tools/goctl@latest")
		if err != nil {
			return err
		}
	}
	if _, err = env.LookPath("goctl"); err != nil {
		return errors.New("goctl is not installed")
	}

	// Install frame-specific tools
	if frameType == "rpc" {
		// goctl env check
		resp, err := execx.Run("goctl env check --install --verbose --force", "")
		if err != nil {
			return fmt.Errorf("goctl env check failed, %s", resp)
		}
		if verbose {
			fmt.Println(resp)
		}

		// chmod +x protoc
		protocPath, err := env.LookPath("protoc")
		if err != nil {
			return err
		}
		err = os.Chmod(protocPath, 0755)
		if err != nil {
			return err
		}

		_, err = env.LookPath("protoc-gen-openapiv2")
		if err != nil {
			fmt.Printf("%s protoc-gen-openapiv2\n", color.WithColor("Installing tool", color.FgGreen))
			if err = golang.Install("github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest"); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc-gen-openapiv2"); err != nil {
			return errors.New("protoc-gen-openapiv2 is not installed")
		}
		if verbose {
			log.Success("\n[jzero-env] protoc-gen-openapiv2 is installed")
		}

		// protoc-gen-doc
		_, err = env.LookPath("protoc-gen-doc")
		if err != nil {
			fmt.Printf("%s protoc-gen-doc\n", color.WithColor("Installing tool", color.FgGreen))
			if err = golang.Install("github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest"); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc-gen-doc"); err != nil {
			return errors.New("protoc-gen-doc is not installed")
		}
		if verbose {
			log.Success("\n[jzero-env] protoc-gen-doc is installed")
		}
	}
	if verbose {
		log.Success("\n[jzero-env] congratulations! your jzero environment is ready!")
	}
	return nil
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: `Check and install all needed tools`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunCheck(true)
		cobra.CheckErr(err)
	},
}

func GetCommand() *cobra.Command {
	return checkCmd
}
