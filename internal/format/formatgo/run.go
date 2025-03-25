package formatgo

import (
	"strings"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/pkg"
	"github.com/jzero-io/jzero/pkg/gitstatus"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
)

/*
	Use gofumpt to format go code
*/

func Run() error {
	check()

	files := getFormatFiles()

	return execFormat(files)
}

func getFormatFiles() []string {
	if config.C.Format.GitChange {
		files, _, err := gitstatus.ChangedFiles(".", ".go")
		if err == nil {
			return files
		} else {
			return []string{"."}
		}
	} else {
		return []string{"."}
	}
}

func execFormat(files []string) error {
	if len(files) > 0 {
		args := []string{"gofumpt", "-l", "-w"}
		args = append(args, files...)
		logx.Debugf("execute command: %s", strings.Join(args, " "))
		err := pkg.Run(strings.Join(args, " "), "")
		if err != nil {
			return err
		}
	}
	return nil
}

func check() {
	log := console.NewColorConsole(true)

	// install goctl
	_, err := env.LookPath("gofumpt")
	if err != nil {
		log.Warning(`[jzero-env]: gofumpt is not found in PATH`)
		if err = golang.Install("go install mvdan.cc/gofumpt@latest"); err != nil {
			log.Fatalln("[jzero-env]: gofumpt is not installed, please install it by running: go install mvdan.cc/gofumpt@latest")
		}
	}
	if _, err = env.LookPath("gofumpt"); err != nil {
		log.Fatalln("[jzero-env]: env check failed, gofumpt is not installed")
	}
}
