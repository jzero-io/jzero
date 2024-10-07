package genapi

import (
	"fmt"
	goformat "go/format"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/gitstatus"
	"github.com/jzero-io/jzero/pkg/osx"
	"github.com/jzero-io/jzero/pkg/templatex"
)

type JzeroApi struct {
	Wd     string
	Module string

	config.GenConfig

	ApiFiles          []string
	GenCodeApiFiles   []string
	ApiSpecMap        map[string]*spec.ApiSpec
	GenCodeApiSpecMap map[string]*spec.ApiSpec
}

func (ja *JzeroApi) Gen() error {
	apiDirName := filepath.Join("desc", "api")

	var allHandlerFiles []HandlerFile
	var allLogicFiles []LogicFile

	if !pathx.FileExists(apiDirName) {
		return nil
	}

	fmt.Printf("%s to generate api code.\n", color.WithColor("Start", color.FgGreen))

	// format api dir
	command := fmt.Sprintf("goctl api format --dir %s", apiDirName)
	_, err := execx.Run(command, ja.Wd)
	if err != nil {
		return err
	}

	apiFiles, err := desc.FindApiFiles(apiDirName)
	if err != nil {
		return err
	}

	ja.ApiFiles = apiFiles
	ja.ApiSpecMap = make(map[string]*spec.ApiSpec, len(apiFiles))
	ja.GenCodeApiSpecMap = make(map[string]*spec.ApiSpec, len(apiFiles))

	for _, v := range apiFiles {
		apiSpec, err := parser.Parse(v, nil)
		if err != nil {
			return err
		}
		ja.ApiSpecMap[v] = apiSpec

		logicFiles, err := ja.getAllLogicFiles(v, apiSpec)
		if err != nil {
			return err
		}
		allLogicFiles = append(allLogicFiles, logicFiles...)

		handlerFiles, err := ja.getAllHandlerFiles(v, apiSpec)
		if err != nil {
			return err
		}
		allHandlerFiles = append(allHandlerFiles, handlerFiles...)
	}

	var genCodeApiFiles []string

	switch {
	case ja.GitDiff && len(ja.Desc) == 0:
		// 从 git status 获取变动的文件生成
		m, _, err := gitstatus.ChangedFiles(ja.ApiGitDiffPath, ".api")
		if err == nil {
			// 获取变动的 api 文件
			genCodeApiFiles = append(genCodeApiFiles, m...)
			for _, file := range m {
				ja.GenCodeApiSpecMap[file] = ja.ApiSpecMap[file]
			}
		}
	case len(ja.Desc) > 0:
		// 从指定的 desc 文件夹或者文件生成
		for _, v := range ja.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".api" {
					genCodeApiFiles = append(genCodeApiFiles, filepath.Join(strings.Split(filepath.ToSlash(v), "/")...))
				}
			} else {
				specifiedApiFiles, err := desc.FindApiFiles(v)
				if err != nil {
					return err
				}
				genCodeApiFiles = append(genCodeApiFiles, specifiedApiFiles...)
				for _, saf := range specifiedApiFiles {
					ja.GenCodeApiSpecMap[saf] = ja.ApiSpecMap[saf]
				}
			}
		}
	default:
		// 否则就是全量的 api 文件
		genCodeApiFiles = ja.ApiFiles
		ja.GenCodeApiSpecMap = ja.ApiSpecMap
	}
	ja.GenCodeApiFiles = genCodeApiFiles

	err = ja.generateApiCode()
	if err != nil {
		return err
	}

	// 处理多余后缀
	if ja.RemoveSuffix {
		for _, file := range allHandlerFiles {
			if _, ok := ja.GenCodeApiSpecMap[file.ApiFilepath]; ok {
				if err = ja.removeHandlerSuffix(file.Path); err != nil {
					return errors.Wrapf(err, "rewrite %s", file.Path)
				}
			}
		}
		for _, file := range allLogicFiles {
			if _, ok := ja.GenCodeApiSpecMap[file.ApiFilepath]; ok {
				if err = ja.removeLogicSuffix(file.Path); err != nil {
					return errors.Wrapf(err, "rewrite %s", file.Path)
				}
			}
		}
	}

	// 自动替换 logic 层的 request 和 response name
	if ja.ChangeLogicTypes {
		for _, file := range allLogicFiles {
			if _, ok := ja.GenCodeApiSpecMap[file.ApiFilepath]; ok {
				if err := ja.changeLogicTypes(file); err != nil {
					console.Warning("[warning]: rewrite %s meet error %v", file.Path, err)
					continue
				}
			}
		}
	}

	// 将 types.go 分 group 或者分 dir
	err = ja.separateTypesGo(allLogicFiles, allHandlerFiles)
	if err != nil {
		return err
	}

	fmt.Println(color.WithColor("Done", color.FgGreen))
	return nil
}

func (ja *JzeroApi) generateApiCode() error {
	for _, file := range ja.GenCodeApiFiles {
		if parse, ok := ja.GenCodeApiSpecMap[file]; ok {
			for _, group := range parse.Service.Groups {
				if ja.RegenApiHandler {
					dirFile, err := os.ReadDir(filepath.Join(ja.Wd, "internal", "handler", group.GetAnnotation("group")))
					if err != nil {
						return err
					}
					for _, v := range dirFile {
						if !v.IsDir() {
							_ = os.Remove(filepath.Join(ja.Wd, "internal", "handler", group.GetAnnotation("group"), v.Name()))
						}
					}
				}
				if ja.SplitApiTypesDir {
					_ = os.RemoveAll(filepath.Join(ja.Wd, "internal", "types", group.GetAnnotation("group")))
				}
			}
		}
	}

	// 处理模板
	var goctlHome string
	if !pathx.FileExists(filepath.Join(config.C.Gen.Home, "go-zero", "api")) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
		err = embeded.WriteTemplateDir(filepath.Join("go-zero", "api"), filepath.Join(tempDir, "api"))
		if err != nil {
			return err
		}
		goctlHome = tempDir
	} else {
		goctlHome = filepath.Join(config.C.Gen.Home, "go-zero")
	}
	logx.Debugf("goctl_home = %s", goctlHome)

	var handlerImports desc.ImportLines
	var allRoutesGoBody string
	var allRoutesGoBodyMap sync.Map

	var eg errgroup.Group
	eg.SetLimit(len(ja.ApiFiles))
	for _, v := range ja.ApiFiles {
		cv := v
		eg.Go(func() error {
			routesGoBody, err := ja.getRoutesGoBody(cv)
			if err != nil {
				return err
			}
			if routesGoBody != "" {
				allRoutesGoBodyMap.Store(cv, routesGoBody)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	for _, v := range ja.ApiFiles {
		if s, ok := allRoutesGoBodyMap.Load(v); ok {
			allRoutesGoBody += cast.ToString(s) + "\n"
		}
	}

	eg.SetLimit(len(ja.GenCodeApiFiles))
	for _, v := range ja.GenCodeApiFiles {
		cv := v
		if len(ja.ApiSpecMap[cv].Service.Routes()) > 0 {
			eg.Go(func() error {
				dir := "."
				fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), cv)
				_ = os.Remove(filepath.Join("internal", "types", "types.go"))
				filename := desc.GetApiFrameEtcFilename(ja.Wd, ja.Style)
				if filename != "etc.yaml" {
					_ = os.Remove(filepath.Join("etc", filename))
				}
				command := fmt.Sprintf("goctl api go --api %s --dir %s --home %s --style %s", cv, dir, goctlHome, ja.Style)
				logx.Debugf("command: %s", command)
				if _, err := execx.Run(command, ja.Wd); err != nil {
					return errors.Wrapf(err, "api file: %s", cv)
				}
				return nil
			})
		}
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	exist := make(map[string]struct{})
	for _, v := range ja.ApiFiles {
		for _, g := range ja.ApiSpecMap[v].Service.Groups {
			if _, ok := exist[g.GetAnnotation("group")]; ok {
				continue
			}
			exist[g.GetAnnotation("group")] = struct{}{}
			if g.GetAnnotation("group") != "" {
				handlerImports = append(handlerImports, fmt.Sprintf(`%s "%s/internal/handler/%s"`, strings.ToLower(strings.ReplaceAll(g.GetAnnotation("group"), "/", "")), ja.Module, g.GetAnnotation("group")))
			}
		}
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Routes":         allRoutesGoBody,
		"Module":         ja.Module,
		"HandlerImports": handlerImports,
	}, embeded.ReadTemplateFile(filepath.Join("app", "internal", "handler", "routes.go.tpl")))
	if err != nil {
		return err
	}
	source, err := goformat.Source(template)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join("internal", "handler", "routes.go"), source, 0o644)
}
