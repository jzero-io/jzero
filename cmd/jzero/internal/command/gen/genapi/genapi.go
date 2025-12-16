package genapi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/format"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"golang.org/x/sync/errgroup"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/filex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/gitstatus"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
)

type JzeroApi struct {
	Module string

	ApiFiles          []string
	GenCodeApiFiles   []string
	ApiSpecMap        map[string]*spec.ApiSpec
	GenCodeApiSpecMap map[string]*spec.ApiSpec
}

type (
	ImportLines []string

	RegisterLines []string
)

func (l ImportLines) String() string {
	return "\n\n\t" + strings.Join(l, "\n\t")
}

func (l RegisterLines) String() string {
	return "\n\t\t" + strings.Join(l, "\n\t\t")
}

func (ja *JzeroApi) Gen() error {
	if !pathx.FileExists(config.C.ApiDir()) {
		return nil
	}

	apiFiles, err := desc.FindApiFiles(config.C.ApiDir())
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
	case config.C.Gen.GitChange && gitstatus.IsGitRepo(filepath.Join(config.C.Wd())) && len(config.C.Gen.Desc) == 0:
		// 从 git status 获取变动的文件生成
		m, _, err := gitstatus.ChangedFiles(config.C.ApiDir(), ".api")
		if err == nil {
			// 获取变动的 api 文件
			genCodeApiFiles = append(genCodeApiFiles, m...)
			for _, file := range m {
				ja.GenCodeApiSpecMap[file] = ja.ApiSpecMap[file]
			}
		}
	case len(config.C.Gen.Desc) > 0:
		// 从指定的 desc 文件夹或者文件生成
		for _, v := range config.C.Gen.Desc {
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
					ja.GenCodeApiSpecMap[filepath.Clean(saf)] = ja.ApiSpecMap[filepath.Clean(saf)]
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
	for _, v := range config.C.Gen.DescIgnore {
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

	if len(ja.GenCodeApiFiles) == 0 {
		return nil
	}

	if !config.C.Quiet {
		fmt.Printf("%s to generate api code from api files\n", console.Green("Start"))
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

	if !config.C.Quiet {
		fmt.Println(console.Green("Done"))
	}
	return nil
}

func (ja *JzeroApi) generateApiCode() error {
	for _, file := range ja.GenCodeApiFiles {
		if parse, ok := ja.GenCodeApiSpecMap[file]; ok {
			for _, group := range parse.Service.Groups {
				dirFile, err := os.ReadDir(filepath.Join(config.C.Wd(), "internal", "handler", group.GetAnnotation("group")))
				if err == nil {
					for _, v := range dirFile {
						if !v.IsDir() {
							_ = os.Remove(filepath.Join(config.C.Wd(), "internal", "handler", group.GetAnnotation("group"), v.Name()))
						}
					}
				}
			}
		}
	}

	// 处理模板
	var goctlHome string
	tempDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 先写入内置模板
	err = embeded.WriteTemplateDir(filepath.Join("go-zero", "api"), filepath.Join(tempDir, "api"))
	if err != nil {
		return err
	}

	// 如果用户自定义了模板，则复制覆盖
	customTemplatePath := filepath.Join(config.C.Home, "go-zero", "api")
	if pathx.FileExists(customTemplatePath) {
		err = filex.CopyDir(customTemplatePath, filepath.Join(tempDir, "api"))
		if err != nil {
			return err
		}
	}

	goctlHome = tempDir
	logx.Debugf("goctl_home = %s", goctlHome)

	var handlerImports ImportLines
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
			logicFiles, err := ja.getAllLogicFiles(v, ja.ApiSpecMap[v])
			if err != nil {
				return err
			}

			handlerFiles, err := ja.getAllHandlerFiles(v, ja.ApiSpecMap[v])
			if err != nil {
				return err
			}

			dir := "."
			if !config.C.Quiet {
				fmt.Printf("%s api file %s\n", console.Green("Using"), v)
			}

			if err = format.ApiFormatByPath(v, false); err != nil {
				return errors.Wrapf(err, "format api file: %s", v)
			}

			command := fmt.Sprintf("goctl api go --api %s --dir %s --home %s --style %s", v, dir, goctlHome, config.C.Style)
			logx.Debugf("command: %s", command)
			if _, err := execx.Run(command, config.C.Wd()); err != nil {
				return errors.Wrapf(err, "api file: %s", v)
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
				handlerImports = append(handlerImports, fmt.Sprintf(`%s "%s/internal/handler/%s"`, strings.ReplaceAll(g.GetAnnotation("group"), "/", ""), ja.Module, g.GetAnnotation("group")))
			}
		}
	}

	template, err := templatex.ParseTemplate(filepath.Join("api", "routes.go.tpl"), map[string]any{
		"Routes":         allRoutesGoBody,
		"Module":         ja.Module,
		"HandlerImports": lo.Uniq(handlerImports),
	}, embeded.ReadTemplateFile(filepath.Join("api", "routes.go.tpl")))
	if err != nil {
		return err
	}
	process, err := gosimports.Process("", template, nil)
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join("internal", "handler", "routes.go"), process, 0o644); err != nil {
		return err
	}

	if config.C.Gen.Route2Code {
		if !config.C.Quiet {
			fmt.Printf("%s to generate internal/handler/route2code.go\n", console.Green("Start"))
		}
		if route2CodeBytes, err := ja.genRoute2Code(); err != nil {
			return err
		} else {
			if err = os.WriteFile(filepath.Join("internal", "handler", "route2code.go"), route2CodeBytes, 0o644); err != nil {
				return err
			}
		}
		if !config.C.Quiet {
			fmt.Printf("%s", console.Green("Done\n"))
		}
	}
	return nil
}
