package gen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/embeded"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

type JzeroApi struct {
	Wd     string
	Module string
	Style  string
}

func (ja *JzeroApi) Gen() error {
	apiDirName := filepath.Join(ja.Wd, "app", "desc", "api")
	if pathx.FileExists(apiDirName) {
		// format api dir
		command := fmt.Sprintf("goctl api format --dir %s", apiDirName)
		_, err := execx.Run(command, ja.Wd)
		if err != nil {
			return err
		}

		fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))
		mainApiFilePath := GetMainApiFilePath(apiDirName)

		err = generateApiCode(ja.Wd, mainApiFilePath, ja.Style)
		if err != nil {
			return err
		}
		// goctl-types. make types.go separate by group
		err = separateTypesGoByGoctlTypesPlugin(ja.Wd, mainApiFilePath, ja.Style)
		if err != nil {
			return err
		}
		_ = os.Remove(mainApiFilePath)
		fmt.Println(color.WithColor("Done", color.FgGreen))
	}

	return nil
}

func getRouteApiFilePath(apiDirName string) []string {
	var apiFiles []string
	_ = filepath.Walk(apiDirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".api" {
			apiSpec, err := parser.Parse(path, nil)
			if err != nil {
				return err
			}
			if len(apiSpec.Service.Routes()) > 0 {
				rel, err := filepath.Rel(apiDirName, path)
				if err != nil {
					return err
				}
				apiFiles = append(apiFiles, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	return apiFiles
}

func generateApiCode(wd string, mainApiFilePath, style string) error {
	if mainApiFilePath == "" {
		return errors.New("empty mainApiFilePath")
	}

	fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), mainApiFilePath)
	command := fmt.Sprintf("goctl api go --api %s --dir ./app --home %s --style %s ", mainApiFilePath, filepath.Join(embeded.Home, "go-zero"), style)
	if _, err := execx.Run(command, wd); err != nil {
		return err
	}
	return nil
}

func separateTypesGoByGoctlTypesPlugin(wd string, mainApiFilePath, style string) error {
	command := fmt.Sprintf("goctl api plugin -plugin goctl-types=\"gen\" -api %s --dir ./app --style %s\n", mainApiFilePath, style)
	if _, err := execx.Run(command, wd); err != nil {
		return err
	}
	return nil
}
