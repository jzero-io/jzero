package genapi

import (
	"fmt"
	goformat "go/format"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
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
			return errors.Wrapf(err, "parse %s", v)
		}
		ja.ApiSpecMap[v] = apiSpec
	}

	var genCodeApiFiles []string

	switch {
	case ja.GitChange && len(ja.Desc) == 0:
		// 从 git status 获取变动的文件生成
		m, _, err := gitstatus.ChangedFiles(ja.ApiGitChangePath, ".api")
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
					genCodeApiFiles = append(genCodeApiFiles, filepath.Join(strings.Split(filepath.ToSlash(filepath.Clean(v)), "/")...))
					ja.GenCodeApiSpecMap[filepath.Clean(v)] = ja.ApiSpecMap[filepath.Clean(v)]
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
		// clone 一份 gen code api spec map
		for k, v := range ja.ApiSpecMap {
			ja.GenCodeApiSpecMap[k] = v
		}
	}
	ja.GenCodeApiFiles = genCodeApiFiles

	// ignore api desc
	for _, v := range ja.DescIgnore {
		if !osx.IsDir(v) {
			if filepath.Ext(v) == ".api" {
				// delete item in genCodeApiFiles by filename
				ja.GenCodeApiFiles = lo.Reject(ja.GenCodeApiFiles, func(item string, _ int) bool {
					return item == v
				})
				// delete map key
				delete(ja.GenCodeApiSpecMap, v)
			}
		} else {
			specifiedApiFiles, err := desc.FindApiFiles(v)
			if err != nil {
				return err
			}
			for _, saf := range specifiedApiFiles {
				ja.GenCodeApiFiles = lo.Reject(ja.GenCodeApiFiles, func(item string, _ int) bool {
					return item == saf
				})
				delete(ja.GenCodeApiSpecMap, saf)
			}
		}
	}

	err = ja.generateApiCode()
	if err != nil {
		return err
	}

	// 将 types.go 分 group 或者分 dir
	err = ja.separateTypesGo()
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
					if err == nil {
						for _, v := range dirFile {
							if !v.IsDir() {
								_ = os.Remove(filepath.Join(ja.Wd, "internal", "handler", group.GetAnnotation("group"), v.Name()))
							}
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

	for _, v := range ja.GenCodeApiFiles {
		if len(ja.ApiSpecMap[v].Service.Routes()) > 0 {
			dir := "."
			fmt.Printf("%s api file %s\n", color.WithColor("Using", color.FgGreen), v)
			command := fmt.Sprintf("goctl api go --api %s --dir %s --home %s --style %s", v, dir, goctlHome, ja.Style)
			logx.Debugf("command: %s", command)
			if _, err := execx.Run(command, ja.Wd); err != nil {
				return errors.Wrapf(err, "api file: %s", v)
			}

			logicFiles, err := ja.getAllLogicFiles(v, ja.ApiSpecMap[v])
			if err != nil {
				return err
			}

			handlerFiles, err := ja.getAllHandlerFiles(v, ja.ApiSpecMap[v])
			if err != nil {
				return err
			}

			// patch handler
			for _, file := range handlerFiles {
				if _, ok := ja.GenCodeApiSpecMap[file.ApiFilepath]; ok {
					if err = ja.patchHandler(file); err != nil {
						return errors.Wrapf(err, "rewrite %s", file.Path)
					}
				}
			}
			for _, file := range logicFiles {
				if _, ok := ja.GenCodeApiSpecMap[file.DescFilepath]; ok {
					if err = ja.patchLogic(file); err != nil {
						return errors.Wrapf(err, "rewrite %s", file.Path)
					}
				}
			}
		}
	}

	for _, v := range ja.ApiFiles {
		for _, g := range ja.ApiSpecMap[v].Service.Groups {
			if g.GetAnnotation("group") != "" {
				handlerImports = append(handlerImports, fmt.Sprintf(`%s "%s/internal/handler/%s"`, strings.ToLower(strings.ReplaceAll(g.GetAnnotation("group"), "/", "")), ja.Module, g.GetAnnotation("group")))
			}
		}
	}

	template, err := templatex.ParseTemplate(map[string]any{
		"Routes":         allRoutesGoBody,
		"Module":         ja.Module,
		"HandlerImports": lo.Uniq(handlerImports),
	}, embeded.ReadTemplateFile(filepath.Join("plugins", "api", "routes.go.tpl")))
	if err != nil {
		return err
	}
	source, err := goformat.Source(template)
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join("internal", "handler", "routes.go"), source, 0o644); err != nil {
		return err
	}

	if ja.Route2Code {
		if route2CodeBytes, err := ja.genRoute2Code(); err != nil {
			return err
		} else {
			if err = os.WriteFile(filepath.Join("internal", "handler", "route2code.go"), route2CodeBytes, 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}
