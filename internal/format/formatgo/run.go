package formatgo

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/pkg"
	"github.com/jzero-io/jzero/pkg/gitstatus"
)

/*
	Use gofumpt to format go code
*/

var rxCodeGenerated = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

func Run() error {
	check()

	files := getFormatFiles()

	files = filterFiles(files)

	return execFormat(files)
}

func getFormatFiles() []string {
	if config.C.Format.GitChange {
		files, _, err := gitstatus.ChangedFiles(".", ".go")
		if err == nil {
			return files
		}
		return []string{"."}
	} else {
		return []string{"."}
	}
}

func filterFiles(files []string) []string {
	var result []string

	mr.ForEach(func(source chan<- string) {
		for _, v := range files {
			source <- v
		}
	}, func(item string) {
		line, _ := readFirstLine(item)
		if !rxCodeGenerated.MatchString(line) {
			result = append(result, item)
		}
	})

	return result
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

func readFirstLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", scanner.Err()
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
