package formatgo

import (
	"bufio"
	"os"
	"regexp"

	"github.com/jzero-io/go_fmt/gofmtapi"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"

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
	logx.Debugf("format files: %v", opt.Files)
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

	var eg errgroup.Group
	eg.SetLimit(len(files))
	for _, v := range files {
		eg.Go(func() error {
			line, _ := readFirstLine(v)
			if !rxCodeGenerated.MatchString(line) {
				result = append(result, v)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil
	}

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
