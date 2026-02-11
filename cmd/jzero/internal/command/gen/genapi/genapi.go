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

func (ja *JzeroApi) Gen() (map[string]*spec.ApiSpec, error) {
	if !pathx.FileExists(config.C.ApiDir()) {
		return nil, nil
	}

	apiFiles, err := desc.FindRouteApiFiles(config.C.ApiDir())
	if err != nil {
		return nil, err
	}

	apiSpecMap := make(map[string]*spec.ApiSpec, len(apiFiles))
	genCodeApiSpecMap := make(map[string]*spec.ApiSpec, len(apiFiles))

	for _, v := range apiFiles {
		apiSpec, err := parser.Parse(v, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "parse %s", v)
		}
		apiSpecMap[v] = apiSpec
	}

	// 收集当前文件的路由（不包含 import）
	currentRoutesMap := make(map[string][]spec.Route, len(apiFiles))
	for _, v := range apiFiles {
		routes, err := desc.ParseCurrentApiRoutes(v)
		if err != nil {
			return nil, errors.Wrapf(err, "parse current routes %s", v)
		}
		currentRoutesMap[v] = routes
	}

	// 记录哪些文件被其他文件 import 了
	importedFiles := make(map[string]bool)
	for _, v := range apiFiles {
		for _, imp := range apiSpecMap[v].Imports {
			importPath := desc.ResolveImportPath(v, imp.Value, config.C.ApiDir())
			if importPath != "" {
				importedFiles[importPath] = true
			}
		}
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
				genCodeApiSpecMap[file] = apiSpecMap[file]
			}
		}
	case len(config.C.Gen.Desc) > 0:
		// 从指定的 desc 文件夹或者文件生成
		for _, v := range config.C.Gen.Desc {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".api" {
					genCodeApiFiles = append(genCodeApiFiles, filepath.Join(strings.Split(filepath.ToSlash(filepath.Clean(v)), "/")...))
					genCodeApiSpecMap[filepath.Clean(v)] = apiSpecMap[filepath.Clean(v)]
				}
			} else {
				specifiedApiFiles, err := desc.FindApiFiles(v)
				if err != nil {
					return nil, err
				}
				genCodeApiFiles = append(genCodeApiFiles, specifiedApiFiles...)
				for _, saf := range specifiedApiFiles {
					genCodeApiSpecMap[filepath.Clean(saf)] = apiSpecMap[filepath.Clean(saf)]
				}
			}
		}
	default:
		// 否则就是全量的 api 文件
		genCodeApiFiles = apiFiles
		// clone 一份 gen code api spec map
		for k, v := range apiSpecMap {
			genCodeApiSpecMap[k] = v
		}
	}

	// ignore api desc
	for _, v := range config.C.Gen.DescIgnore {
		if !osx.IsDir(v) {
			if filepath.Ext(v) == ".api" {
				// delete item in genCodeApiFiles by filename
				genCodeApiFiles = lo.Reject(genCodeApiFiles, func(item string, _ int) bool {
					return item == filepath.Clean(v)
				})
				apiFiles = lo.Reject(apiFiles, func(item string, _ int) bool {
					return item == filepath.Clean(v)
				})
				// delete map key
				delete(genCodeApiSpecMap, filepath.Clean(v))
				delete(apiSpecMap, filepath.Clean(v))
			}
		} else {
			specifiedApiFiles, err := desc.FindApiFiles(v)
			if err != nil {
				return nil, err
			}
			for _, saf := range specifiedApiFiles {
				genCodeApiFiles = lo.Reject(genCodeApiFiles, func(item string, _ int) bool {
					return item == saf
				})
				apiFiles = lo.Reject(apiFiles, func(item string, _ int) bool {
					return item == saf
				})
				delete(genCodeApiSpecMap, saf)
				delete(apiSpecMap, saf)
			}
		}
	}

	if len(genCodeApiFiles) == 0 {
		return apiSpecMap, nil
	}

	if !config.C.Quiet {
		fmt.Printf("%s to generate api code from api files\n", console.Green("Start"))
	}

	err = ja.generateApiCode(apiFiles, apiSpecMap, genCodeApiFiles, genCodeApiSpecMap, currentRoutesMap, importedFiles)
	if err != nil {
		return nil, err
	}

	// 将 types.go 分 group 或者分 dir
	err = ja.separateTypesGo(apiFiles, apiSpecMap)
	if err != nil {
		return nil, err
	}

	if !config.C.Quiet {
		fmt.Println(console.Green("Done"))
	}
	return apiSpecMap, nil
}

func (ja *JzeroApi) generateApiCode(apiFiles []string, apiSpecMap map[string]*spec.ApiSpec, genCodeApiFiles []string, genCodeApiSpecMap map[string]*spec.ApiSpec, currentRoutesMap map[string][]spec.Route, importedFiles map[string]bool) error {
	if err := ja.cleanHandlersDir(genCodeApiFiles, genCodeApiSpecMap); err != nil {
		return err
	}

	templateDir, err := ja.prepareTemplateDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(templateDir)

	allRoutesGoBody, err := ja.collectRoutesGoBody(apiFiles, apiSpecMap, currentRoutesMap, importedFiles)
	if err != nil {
		return err
	}

	if err := ja.generateCodeForApiFiles(genCodeApiFiles, apiSpecMap, importedFiles, templateDir); err != nil {
		return err
	}

	if err := ja.patchHandlerAndLogicFiles(genCodeApiFiles, apiSpecMap, genCodeApiSpecMap); err != nil {
		return err
	}

	if err := ja.generateRoutesGoFile(apiFiles, apiSpecMap, importedFiles, allRoutesGoBody); err != nil {
		return err
	}

	if config.C.Gen.Route2Code {
		if err := ja.generateRoute2CodeFile(apiSpecMap, currentRoutesMap, importedFiles); err != nil {
			return err
		}
	}

	return nil
}

// cleanHandlersDir 清理 handler 目录下的旧文件
func (ja *JzeroApi) cleanHandlersDir(genCodeApiFiles []string, genCodeApiSpecMap map[string]*spec.ApiSpec) error {
	var eg errgroup.Group
	for _, file := range genCodeApiFiles {
		apiSpec, ok := genCodeApiSpecMap[file]
		if !ok {
			continue
		}

		for _, group := range apiSpec.Service.Groups {
			groupAnnotation := group.GetAnnotation("group")
			if groupAnnotation == "" {
				continue
			}

			handlerDir := filepath.Join(config.C.Wd(), "internal", "handler", groupAnnotation)
			eg.Go(func() error {
				dirEntries, err := os.ReadDir(handlerDir)
				if err != nil {
					return nil // 目录不存在或无法读取，忽略错误
				}
				for _, entry := range dirEntries {
					if !entry.IsDir() {
						_ = os.Remove(filepath.Join(handlerDir, entry.Name()))
					}
				}
				return nil
			})
		}
	}
	return eg.Wait()
}

// prepareTemplateDir 准备模板目录，返回临时目录路径
func (ja *JzeroApi) prepareTemplateDir() (string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return "", err
	}

	// 写入内置模板
	if err := embeded.WriteTemplateDir(filepath.Join("go-zero", "api"), filepath.Join(tempDir, "api")); err != nil {
		_ = os.RemoveAll(tempDir)
		return "", err
	}

	// 如果用户自定义了模板，则复制覆盖
	customTemplatePath := filepath.Join(config.C.Home, "go-zero", "api")
	if pathx.FileExists(customTemplatePath) {
		if err := filex.CopyDir(customTemplatePath, filepath.Join(tempDir, "api")); err != nil {
			_ = os.RemoveAll(tempDir)
			return "", err
		}
	}

	logx.Debugf("goctl_home = %s", tempDir)
	return tempDir, nil
}

// collectRoutesGoBody 并发收集所有文件的 routesGoBody
func (ja *JzeroApi) collectRoutesGoBody(apiFiles []string, apiSpecMap map[string]*spec.ApiSpec, currentRoutesMap map[string][]spec.Route, importedFiles map[string]bool) (string, error) {
	var allRoutesGoBodyMap sync.Map

	var eg errgroup.Group
	eg.SetLimit(len(apiFiles))

	for _, apiFile := range apiFiles {
		if importedFiles[apiFile] {
			continue
		}

		currentFile := apiFile
		eg.Go(func() error {
			routesGoBody, err := ja.getRoutesGoBody(currentFile, apiSpecMap, currentRoutesMap)
			if err != nil {
				return err
			}
			if routesGoBody != "" {
				allRoutesGoBodyMap.Store(currentFile, routesGoBody)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return "", err
	}

	var allRoutesGoBody strings.Builder
	for _, apiFile := range apiFiles {
		if body, ok := allRoutesGoBodyMap.Load(apiFile); ok {
			allRoutesGoBody.WriteString(cast.ToString(body))
			allRoutesGoBody.WriteString("\n")
		}
	}

	return allRoutesGoBody.String(), nil
}

// generateCodeForApiFiles 为所有 API 文件生成代码
func (ja *JzeroApi) generateCodeForApiFiles(genCodeApiFiles []string, apiSpecMap map[string]*spec.ApiSpec, importedFiles map[string]bool, templateDir string) error {
	// 按 group 分组，同一 group 的文件串行处理，不同 group 并发处理
	groupToFiles := make(map[string][]string)
	for _, apiFile := range genCodeApiFiles {
		if importedFiles[apiFile] {
			continue
		}
		if len(apiSpecMap[apiFile].Service.Routes()) == 0 {
			continue
		}

		// 收集该文件的所有 group
		groups := make(map[string]struct{})
		for _, g := range apiSpecMap[apiFile].Service.Groups {
			if groupAnnotation := g.GetAnnotation("group"); groupAnnotation != "" {
				groups[groupAnnotation] = struct{}{}
			}
		}

		// 如果没有 group，使用默认分组
		if len(groups) == 0 {
			groupToFiles[""] = append(groupToFiles[""], apiFile)
		} else {
			for group := range groups {
				groupToFiles[group] = append(groupToFiles[group], apiFile)
			}
		}
	}

	// 并发处理不同 group
	var eg errgroup.Group
	for _, files := range groupToFiles {
		currentFiles := files
		eg.Go(func() error {
			// 同一 group 内的文件串行处理
			for _, apiFile := range currentFiles {
				if !config.C.Quiet {
					fmt.Printf("%s api file %s \n", console.Green("Using"), apiFile)
				}

				if err := format.ApiFormatByPath(apiFile, false); err != nil {
					return errors.Wrapf(err, "format api file: %s", apiFile)
				}

				command := fmt.Sprintf("goctl api go --api %s --dir %s --home %s --style %s", apiFile, ".", templateDir, config.C.Style)
				logx.Debugf("command: %s", command)

				if _, err := execx.Run(command, config.C.Wd()); err != nil {
					return errors.Wrapf(err, "api file: %s", apiFile)
				}
			}
			return nil
		})
	}

	return eg.Wait()
}

// patchHandlerAndLogicFiles 并发 patch handler 和 logic 文件
func (ja *JzeroApi) patchHandlerAndLogicFiles(genCodeApiFiles []string, apiSpecMap map[string]*spec.ApiSpec, genCodeApiSpecMap map[string]*spec.ApiSpec) error {
	var eg errgroup.Group

	for _, apiFile := range genCodeApiFiles {
		if len(apiSpecMap[apiFile].Service.Routes()) == 0 {
			continue
		}

		currentFile := apiFile

		eg.Go(func() error {
			logicFiles, err := ja.getAllLogicFiles(currentFile, apiSpecMap[currentFile])
			if err != nil {
				return err
			}

			handlerFiles, err := ja.getAllHandlerFiles(currentFile, apiSpecMap[currentFile])
			if err != nil {
				return err
			}

			// Patch handler files
			for _, file := range handlerFiles {
				if _, ok := genCodeApiSpecMap[file.ApiFilepath]; ok {
					if err = ja.patchHandler(file, genCodeApiSpecMap); err != nil {
						return errors.Wrapf(err, "rewrite %s", file.Path)
					}
				}
			}

			// Patch logic files
			for _, file := range logicFiles {
				if _, ok := genCodeApiSpecMap[file.DescFilepath]; ok {
					if err = ja.patchLogic(file, genCodeApiSpecMap); err != nil {
						return errors.Wrapf(err, "rewrite %s", file.Path)
					}
				}
			}

			return nil
		})
	}

	return eg.Wait()
}

// generateRoutesGoFile 生成 routes.go 文件
func (ja *JzeroApi) generateRoutesGoFile(apiFiles []string, apiSpecMap map[string]*spec.ApiSpec, importedFiles map[string]bool, allRoutesGoBody string) error {
	var handlerImports ImportLines

	for _, apiFile := range apiFiles {
		if importedFiles[apiFile] {
			continue
		}

		for _, group := range apiSpecMap[apiFile].Service.Groups {
			groupAnnotation := group.GetAnnotation("group")
			if groupAnnotation != "" {
				importPath := fmt.Sprintf(`%s "%s/internal/handler/%s"`,
					strings.ReplaceAll(groupAnnotation, "/", ""), ja.Module, groupAnnotation)
				handlerImports = append(handlerImports, importPath)
			}
		}
	}

	templateContent, err := templatex.ParseTemplate(
		filepath.Join("api", "routes.go.tpl"),
		map[string]any{
			"Routes":         allRoutesGoBody,
			"Module":         ja.Module,
			"HandlerImports": lo.Uniq(handlerImports),
		},
		embeded.ReadTemplateFile(filepath.Join("api", "routes.go.tpl")),
	)
	if err != nil {
		return err
	}

	process, err := gosimports.Process("", templateContent, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("internal", "handler", "routes.go"), process, 0o644)
}

// generateRoute2CodeFile 生成 route2code.go 文件
func (ja *JzeroApi) generateRoute2CodeFile(apiSpecMap map[string]*spec.ApiSpec, currentRoutesMap map[string][]spec.Route, importedFiles map[string]bool) error {
	if !config.C.Quiet {
		fmt.Printf("%s to generate internal/handler/route2code.go\n", console.Green("Start"))
	}

	route2CodeBytes, err := ja.genRoute2Code(apiSpecMap, currentRoutesMap, importedFiles)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join("internal", "handler", "route2code.go"), route2CodeBytes, 0o644); err != nil {
		return err
	}

	if !config.C.Quiet {
		fmt.Printf("%s", console.Green("Done\n"))
	}

	return nil
}
