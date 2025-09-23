/*
Copyright © 2024 jaronnie <jaron@jaronnie.com>

*/

package check

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/downloader"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
	"github.com/zeromicro/go-zero/tools/goctl/util/zipx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
)

var toolVersionCheck = map[string]string{
	"protoc":               "32.0",
	"goctl":                "1.9.0",
	"protoc-gen-go":        "1.36.8",
	"protoc-gen-go-grpc":   "1.5.1",
	"protoc-gen-openapiv2": "2.27.2",
	"protoc-gen-doc":       "1.5.1",
}

// RunCheck executes the check logic and can be called from other places
func RunCheck(all bool) error {
	// Detect frame type
	frameType, err := desc.GetFrameType()
	if err != nil {
		return err
	}
	if frameType == "" && !all {
		return nil
	}

	// install goctl
	_, err = env.LookPath("goctl")
	if err != nil {
		fmt.Printf("%s goctl %s\n", color.WithColor("Installing tool", color.FgGreen), toolVersionCheck["goctl"])
		err = golang.Install(fmt.Sprintf("github.com/zeromicro/go-zero/tools/goctl@v%s", toolVersionCheck["goctl"]))
		if err != nil {
			return err
		}
	}
	if _, err = env.LookPath("goctl"); err != nil {
		return errors.New("goctl is not installed")
	}

	// check goctl version
	goctlVersion := config.C.GoctlVersion()
	checkGoctlVersion, err := version.NewVersion(toolVersionCheck["goctl"])
	if err != nil {
		return err
	}
	if goctlVersion == nil || goctlVersion.LessThan(checkGoctlVersion) {
		fmt.Printf("%s goctl to %s\n", color.WithColor("Upgrading tool", color.FgGreen), toolVersionCheck["goctl"])
		err = golang.Install(fmt.Sprintf("github.com/zeromicro/go-zero/tools/goctl@v%s", toolVersionCheck["goctl"]))
		if err != nil {
			return err
		}
	}

	// Install frame-specific tools
	if frameType == "rpc" || frameType == "gateway" || all {
		// protoc
		_, err = env.LookPath("protoc")
		if err != nil {
			fmt.Printf("%s protoc %s\n", color.WithColor("Installing tool", color.FgGreen), toolVersionCheck["protoc"])
			if err = installProtoc(); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc"); err != nil {
			return errors.New("protoc is not installed")
		}

		// check protoc version
		protocVersion := config.C.ProtocVersion()
		checkProtocVersion, err := version.NewVersion(toolVersionCheck["protoc"])
		if err != nil {
			return err
		}
		if protocVersion == nil || protocVersion.LessThan(checkProtocVersion) {
			fmt.Printf("%s protoc to %s\n", color.WithColor("Upgrading tool", color.FgGreen), toolVersionCheck["protoc"])
			if err = installProtoc(); err != nil {
				return err
			}
		}

		// protoc-gen-go
		_, err = env.LookPath("protoc-gen-go")
		if err != nil {
			fmt.Printf("%s protoc-gen-go %s\n", color.WithColor("Installing tool", color.FgGreen), toolVersionCheck["protoc-gen-go"])
			if err = golang.Install(fmt.Sprintf("google.golang.org/protobuf/cmd/protoc-gen-go@v%s", toolVersionCheck["protoc-gen-go"])); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc-gen-go"); err != nil {
			return errors.New("protoc-gen-go is not installed")
		}

		// check protoc-gen-go version
		protocGenGoVersion := config.C.ProtocGenGoVersion()
		checkProtocGenGoVersion, err := version.NewVersion(toolVersionCheck["protoc-gen-go"])
		if err != nil {
			return err
		}
		if protocGenGoVersion == nil || protocGenGoVersion.LessThan(checkProtocGenGoVersion) {
			fmt.Printf("%s protoc-gen-go to %s\n", color.WithColor("Upgrading tool", color.FgGreen), toolVersionCheck["protoc-gen-go"])
			if err = golang.Install(fmt.Sprintf("google.golang.org/protobuf/cmd/protoc-gen-go@v%s", toolVersionCheck["protoc-gen-go"])); err != nil {
				return err
			}
		}

		// protoc-gen-go-grpc
		_, err = env.LookPath("protoc-gen-go-grpc")
		if err != nil {
			fmt.Printf("%s protoc-gen-go-grpc %s\n", color.WithColor("Installing tool", color.FgGreen), toolVersionCheck["protoc-gen-go-grpc"])
			if err = golang.Install(fmt.Sprintf("google.golang.org/grpc/cmd/protoc-gen-go-grpc@v%s", toolVersionCheck["protoc-gen-go-grpc"])); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc-gen-go-grpc"); err != nil {
			return errors.New("protoc-gen-go-grpc is not installed")
		}

		// check protoc-gen-go-grpc version
		protocGenGoGrpcVersion := config.C.ProtocGenGoGrpcVersion()
		checkProtocGenGoGrpcVersion, err := version.NewVersion(toolVersionCheck["protoc-gen-go-grpc"])
		if err != nil {
			return err
		}
		if protocGenGoGrpcVersion == nil || protocGenGoGrpcVersion.LessThan(checkProtocGenGoGrpcVersion) {
			fmt.Printf("%s protoc-gen-go-grpc to %s\n", color.WithColor("Upgrading tool", color.FgGreen), toolVersionCheck["protoc-gen-go-grpc"])
			if err = golang.Install(fmt.Sprintf("google.golang.org/grpc/cmd/protoc-gen-go-grpc@v%s", toolVersionCheck["protoc-gen-go-grpc"])); err != nil {
				return err
			}
		}

		_, err = env.LookPath("protoc-gen-openapiv2")
		if err != nil {
			fmt.Printf("%s protoc-gen-openapiv2 %s\n", color.WithColor("Installing tool", color.FgGreen), toolVersionCheck["protoc-gen-openapiv2"])
			if err = golang.Install(fmt.Sprintf("github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v%s", toolVersionCheck["protoc-gen-openapiv2"])); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc-gen-openapiv2"); err != nil {
			return errors.New("protoc-gen-openapiv2 is not installed")
		}

		// check protoc-gen-openapiv2 version
		protocGenOpenapiv2Version := config.C.ProtocGenOpenapiv2Version()
		checkProtocGenOpenapiv2Version, err := version.NewVersion(toolVersionCheck["protoc-gen-openapiv2"])
		if err != nil {
			return err
		}
		if protocGenOpenapiv2Version == nil || protocGenOpenapiv2Version.LessThan(checkProtocGenOpenapiv2Version) {
			fmt.Printf("%s protoc-gen-openapiv2 to %s\n", color.WithColor("Upgrading tool", color.FgGreen), toolVersionCheck["protoc-gen-openapiv2"])
			if err = golang.Install(fmt.Sprintf("github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v%s", toolVersionCheck["protoc-gen-openapiv2"])); err != nil {
				return err
			}
		}

		// protoc-gen-doc
		_, err = env.LookPath("protoc-gen-doc")
		if err != nil {
			fmt.Printf("%s protoc-gen-doc %s\n", color.WithColor("Installing tool", color.FgGreen), toolVersionCheck["protoc-gen-doc"])
			if err = golang.Install(fmt.Sprintf("github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v%s", toolVersionCheck["protoc-gen-doc"])); err != nil {
				return err
			}
		}
		if _, err = env.LookPath("protoc-gen-doc"); err != nil {
			return errors.New("protoc-gen-doc is not installed")
		}

		// check protoc-gen-doc version
		protocGenDocVersion := config.C.ProtocGenDocVersion()
		checkProtocGenDocVersion, err := version.NewVersion(toolVersionCheck["protoc-gen-doc"])
		if err != nil {
			return err
		}
		if protocGenDocVersion == nil || protocGenDocVersion.LessThan(checkProtocGenDocVersion) {
			fmt.Printf("%s protoc-gen-doc to %s\n", color.WithColor("Upgrading tool", color.FgGreen), toolVersionCheck["protoc-gen-doc"])
			if err = golang.Install(fmt.Sprintf("github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v%s", toolVersionCheck["protoc-gen-doc"])); err != nil {
				return err
			}
		}
	}
	return nil
}

var url = map[string]string{
	"linux_amd64":   "https://github.com/protocolbuffers/protobuf/releases/download/v32.0/protoc-32.0-linux-x86_64.zip",
	"linux_arm64":   "https://github.com/protocolbuffers/protobuf/releases/download/v32.0/protoc-32.0-linux-aarch_64.zip",
	"darwin_amd64":  "https://github.com/protocolbuffers/protobuf/releases/download/v32.0/protoc-32.0-osx-x86_64.zip",
	"darwin_arm64":  "https://github.com/protocolbuffers/protobuf/releases/download/v32.0/protoc-32.0-osx-aarch_64.zip",
	"windows_amd64": "https://github.com/protocolbuffers/protobuf/releases/download/v32.0/protoc-32.0-win64.zip",
	"windows_arm64": "https://github.com/protocolbuffers/protobuf/releases/download/v32.0/protoc-32.0-win64.zip",
}

const (
	ProtocName  = "protoc"
	ZipFileName = ProtocName + ".zip"
)

func installProtoc() error {
	tempFile := filepath.Join(os.TempDir(), ZipFileName)
	var downloadUrl string
	downloadUrl = url[fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)]
	if downloadUrl == "" {
		return errors.Errorf("not support platform %s_%s", runtime.GOOS, runtime.GOARCH)
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
	err = os.Chmod(protocPath, 0o755)
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
		err := RunCheck(true)
		cobra.CheckErr(err)
	},
}

func GetCommand() *cobra.Command {
	return checkCmd
}
