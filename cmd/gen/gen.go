package gen

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jzero-io/jzero/app/pkg/mod"
	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

var (
	WorkingDir string

	Version string
)

type (
	ImportLines   []string
	RegisterLines []string
)

type ApiFileTypes struct {
	Filepath string
	ApiSpec  spec.ApiSpec
	GenTypes []spec.Type

	Base bool
}

func (l ImportLines) String() string {
	return "\n\n\t" + strings.Join(l, "\n\t")
}

func (l RegisterLines) String() string {
	return "\n\t\t" + strings.Join(l, "\n\t\t")
}

func Gen(_ *cobra.Command, _ []string) error {
	// change dir
	if WorkingDir != "" {
		err := os.Chdir(WorkingDir)
		cobra.CheckErr(err)
	}

	wd, err := os.Getwd()
	cobra.CheckErr(err)
	fmt.Printf("%s working dir %s\n", color.WithColor("Enter", color.FgGreen), wd)

	if embeded.Home == "" {
		home, _ := os.UserHomeDir()
		embeded.Home = filepath.Join(home, ".jzero", Version)
	}

	moduleStruct, err := mod.GetGoMod(wd)
	cobra.CheckErr(err)

	// 正常删除无用文件夹
	defer func() {
		removeExtraFiles(wd)
		os.Exit(0)
	}()

	// 异常删除无用文件夹
	go extraFileHandler(wd)

	// gen rpc code
	jzeroRpc := JzeroRpc{Wd: wd, Module: moduleStruct.Path}
	err = jzeroRpc.Gen()
	cobra.CheckErr(err)

	// 生成 api 代码
	jzeroApi := JzeroApi{Wd: wd, Module: moduleStruct.Path}
	err = jzeroApi.Gen()
	cobra.CheckErr(err)

	// 检测是否包含 sql
	jzeroSql := JzeroSql{Wd: wd}
	err = jzeroSql.Gen()
	cobra.CheckErr(err)

	return nil
}

func extraFileHandler(wd string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			removeExtraFiles(wd)
			os.Exit(-1)
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func removeExtraFiles(wd string) {
	_ = os.RemoveAll(filepath.Join(wd, "app", "etc"))
	_ = os.Remove(filepath.Join(wd, "app", fmt.Sprintf("%s.go", GetApiServiceName(filepath.Join(wd, "app", "desc", "api")))))
}

func init() {
	logx.Disable()
}
