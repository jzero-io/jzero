package gen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/embeded"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

type JzeroApi struct {
	Wd           string
	Module       string
	Style        string
	RemoveSuffix bool
}

type HandlerFile struct {
	Path string
	Skip bool
}

type LogicFile struct {
	Path string
	Skip bool
}

func (ja *JzeroApi) Gen() error {
	apiDirName := filepath.Join(ja.Wd, "app", "desc", "api")

	var apiSpec *spec.ApiSpec
	// 实验性功能
	var allHandlerFiles []HandlerFile
	var allLogicFiles []LogicFile

	if pathx.FileExists(apiDirName) {
		// format api dir
		command := fmt.Sprintf("goctl api format --dir %s", apiDirName)
		_, err := execx.Run(command, ja.Wd)
		if err != nil {
			return err
		}

		fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))
		mainApiFilePath := GetMainApiFilePath(apiDirName)
		apiSpec, err = parser.Parse(mainApiFilePath, nil)
		if err != nil {
			return err
		}

		allLogicFiles, err = ja.getAllLogicFiles(apiSpec)
		if err != nil {
			return err
		}

		allHandlerFiles, err = ja.getAllHandlerFiles(apiSpec)
		if err != nil {
			return err
		}

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

	if ja.RemoveSuffix && apiSpec != nil {
		for _, file := range allHandlerFiles {
			if !file.Skip {
				if err := rewriteHandlerGo(file.Path); err != nil {
					return err
				}
			}
		}
		for _, file := range allLogicFiles {
			if !file.Skip {
				if err := rewriteLogicGo(file.Path); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (ja *JzeroApi) getAllHandlerFiles(apiSpec *spec.ApiSpec) ([]HandlerFile, error) {
	var handlerFiles []HandlerFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			formatContent := strings.TrimSuffix(route.Handler, "Handler") + "Handler"
			namingFormat, err := format.FileNamingFormat(ja.Style, formatContent)
			if err != nil {
				return nil, err
			}
			fp := filepath.Join(ja.Wd, "app", "internal", "handler", group.GetAnnotation("group"), namingFormat+".go")

			f := HandlerFile{
				Path: fp,
			}

			if pathx.FileExists(fp) {
				f.Skip = true
			}

			handlerFiles = append(handlerFiles, f)
		}
	}
	return handlerFiles, nil
}

func (ja *JzeroApi) getAllLogicFiles(apiSpec *spec.ApiSpec) ([]LogicFile, error) {
	var handlerFiles []LogicFile
	for _, group := range apiSpec.Service.Groups {
		for _, route := range group.Routes {
			namingFormat, err := format.FileNamingFormat(ja.Style, strings.TrimSuffix(route.Handler, "Handler")+"Logic")
			if err != nil {
				return nil, err
			}

			fp := filepath.Join(ja.Wd, "app", "internal", "logic", group.GetAnnotation("group"), namingFormat+".go")

			f := LogicFile{
				Path: fp,
			}

			if pathx.FileExists(fp) {
				f.Skip = true
			}

			handlerFiles = append(handlerFiles, f)
		}
	}
	return handlerFiles, nil
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

func rewriteHandlerGo(fp string) error {
	return nil
}

func rewriteLogicGo(fp string) error {
	return nil
}
