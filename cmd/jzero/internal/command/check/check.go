/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package check

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/downloader"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
	"github.com/zeromicro/go-zero/tools/goctl/util/zipx"
	"github.com/zeromicro/go-zero/tools/goctl/vars"

	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
)

// RunCheck executes the check logic and can be called from other places
func RunCheck() error {
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
		// protoc
		_, err = env.LookPath("protoc")
		if err != nil {
			fmt.Printf("%s protoc\n", color.WithColor("Installing tool", color.FgGreen))
			if err = installProtoc(); err != nil {
				return err
			}
		}

		// protoc-gen-go
		_, err = env.LookPath("protoc-gen-go")
		if err != nil {
			fmt.Printf("%s protoc-gen-go\n", color.WithColor("Installing tool", color.FgGreen))
			if err = golang.Install("google.golang.org/protobuf/cmd/protoc-gen-go@latest"); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc-gen-go"); err != nil {
			return errors.New("protoc-gen-go is not installed")
		}

		// protoc-gen-go-grpc
		_, err = env.LookPath("protoc-gen-go-grpc")
		if err != nil {
			fmt.Printf("%s protoc-gen-go-grpc\n", color.WithColor("Installing tool", color.FgGreen))
			if err = golang.Install("google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc-gen-go-grpc"); err != nil {
			return errors.New("protoc-gen-go-grpc is not installed")
		}

		_, err = env.LookPath("protoc-gen-openapiv2")
		if err != nil {
			fmt.Printf("%s protoc-gen-openapiv2\n", color.WithColor("Installing tool", color.FgGreen))
			if err = golang.Install("github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest"); err != nil {
				return err
			}
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		}
		if _, err = env.LookPath("protoc-gen-openapiv2"); err != nil {
			return errors.New("protoc-gen-openapiv2 is not installed")
		}

		// protoc-gen-doc
		_, err = env.LookPath("protoc-gen-doc")
		if err != nil {
			fmt.Printf("%s protoc-gen-doc\n", color.WithColor("Installing tool", color.FgGreen))
			if err = golang.Install("github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest"); err != nil {
				return err
			}
			fmt.Printf("%s\n", color.WithColor("Done", color.FgGreen))
		}
		if _, err = env.LookPath("protoc-gen-doc"); err != nil {
			return errors.New("protoc-gen-doc is not installed")
		}
		fmt.Printf("%s\n", color.WithColor("Check done", color.FgGreen))
	}
	return nil
}

var url = map[string]string{
	"linux_32":   "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_32.zip",
	"linux_64":   "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-linux-x86_64.zip",
	"darwin":     "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-osx-x86_64.zip",
	"windows_32": "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-win32.zip",
	"windows_64": "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.4/protoc-3.19.4-win64.zip",
}

const (
	ProtocName  = "protoc"
	ZipFileName = ProtocName + ".zip"
)

func installProtoc() error {
	goos := runtime.GOOS
	tempFile := filepath.Join(os.TempDir(), ZipFileName)
	bit := 32 << (^uint(0) >> 63)
	var downloadUrl string
	switch goos {
	case vars.OsMac:
		downloadUrl = url[vars.OsMac]
	case vars.OsWindows:
		downloadUrl = url[fmt.Sprintf("%s_%d", vars.OsWindows, bit)]
	case vars.OsLinux:
		downloadUrl = url[fmt.Sprintf("%s_%d", vars.OsLinux, bit)]
	default:
		return fmt.Errorf("unsupport OS: %q", goos)
	}

	err := downloader.Download(downloadUrl, tempFile)
	if err != nil {
		return err
	}

	goBin := golang.GoBin()
	protocPath := filepath.Join(goBin, "protoc")
	if runtime.GOOS == "windows" {
		protocPath += ".exe"
	}
	err = zipx.Unpacking(tempFile, goBin, func(f *zip.File) bool {
		return filepath.Base(f.Name) == ProtocName
	})
	if err != nil {
		return err
	}

	// 增加可执行权限
	// chmod +x protoc
	err = os.Chmod(protocPath, 0755)
	if err != nil {
		return err
	}
	return nil
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: `Check and install all needed tools`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunCheck()
		cobra.CheckErr(err)
	},
}

func GetCommand() *cobra.Command {
	return checkCmd
}
