package formatgo

import (
	"bufio"
	"os"
	"regexp"

	"github.com/fsgo/go_fmt/gofmtapi"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/gitstatus"
)

var rxCodeGenerated = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

func Run() error {
	files := getFormatFiles()

	files = filterFiles(files)

	return FormatFiles(files)
}

func FormatFiles(files []string) error {
	gf := gofmtapi.NewFormatter()
	opt := gofmtapi.NewOptions()
	opt.BindFlags()

	opt.DisplayDiff = config.C.Format.DisplayDiff
	opt.Files = files

	if len(opt.Files) == 0 {
		return nil
	}
	return gf.Execute(opt)
}

func getFormatFiles() []string {
	if config.C.Format.GitChange {
		files, _, err := gitstatus.ChangedFiles(".", ".go")
		if err == nil {
			return files
		}
		return []string{"."}
	}
	return []string{"."}
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
