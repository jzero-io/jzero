package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/jparser"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/vars"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/new"
	gconfig "github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/embeded"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/osx"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/templatex"
	"github.com/jzero-io/jzero/cmd/jzero/internal/plugin"
)

func init() {
	Register("go", func(config config.Config) (Generator, error) {
		return &Golang{
			config: &config,
		}, nil
	})
}

type Golang struct {
	config *config.Config

	wd string
}

func (g *Golang) Gen() ([]*GeneratedFile, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	g.wd = wd

	// parse api
	var apiSpecs []*spec.ApiSpec

	if pathx.FileExists(gconfig.C.ApiDir()) {
		var files []string

		switch {
		case len(gconfig.C.Gen.Sdk.Desc) > 0:
			for _, v := range gconfig.C.Gen.Sdk.Desc {
				if !osx.IsDir(v) {
					if filepath.Ext(v) == ".api" {
						files = append(files, v)
					}
				} else {
					specifiedApiFiles, err := jzerodesc.FindApiFiles(v)
					if err != nil {
						return nil, err
					}
					files = append(files, specifiedApiFiles...)
				}
			}
		default:
			files, err = jzerodesc.FindRouteApiFiles(gconfig.C.ApiDir())
			if err != nil {
				return nil, err
			}
		}

		for _, v := range gconfig.C.Gen.Sdk.DescIgnore {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".api" {
					files = lo.Reject(files, func(item string, _ int) bool {
						return item == v
					})
				}
			} else {
				specifiedApiFiles, err := jzerodesc.FindApiFiles(v)
				if err != nil {
					return nil, err
				}
				for _, saf := range specifiedApiFiles {
					files = lo.Reject(files, func(item string, _ int) bool {
						return item == saf
					})
				}
			}
		}

		for _, v := range files {
			apiSpec, err := apiparser.Parse(v, nil)
			if err != nil {
				return nil, errors.Errorf("failed to parse api spec '%s': %v", v, err)
			}
			apiSpecs = append(apiSpecs, apiSpec)
		}
	}

	var protoFiles []string

	if pathx.FileExists(gconfig.C.ProtoDir()) {
		switch {
		case len(gconfig.C.Gen.Sdk.Desc) > 0:
			for _, v := range gconfig.C.Gen.Sdk.Desc {
				if !osx.IsDir(v) {
					if filepath.Ext(v) == ".proto" {
						protoFiles = append(protoFiles, v)
					}
				} else {
					specifiedProtoFiles, err := jzerodesc.GetProtoFilepath(v)
					if err != nil {
						return nil, err
					}
					protoFiles = append(protoFiles, specifiedProtoFiles...)
				}
			}
		default:
			protoFiles, err = jzerodesc.GetProtoFilepath(gconfig.C.ProtoDir())
			if err != nil {
				return nil, err
			}
		}

		for _, v := range gconfig.C.Gen.Sdk.DescIgnore {
			if !osx.IsDir(v) {
				if filepath.Ext(v) == ".proto" {
					protoFiles = lo.Reject(protoFiles, func(item string, _ int) bool {
						return item == v
					})
				}
			} else {
				specifiedProtoFiles, err := jzerodesc.GetProtoFilepath(v)
				if err != nil {
					return nil, err
				}
				for _, saf := range specifiedProtoFiles {
					protoFiles = lo.Reject(protoFiles, func(item string, _ int) bool {
						return item == saf
					})
				}
			}
		}
	}

	var fds []*desc.FileDescriptor

	// parse proto
	var protoParser protoparse.Parser
	if len(protoFiles) > 0 {
		protoParser.ImportPaths = []string{gconfig.C.ProtoDir(), filepath.Join(gconfig.C.ProtoDir(), "third_party")}
		var protoRelFiles []string
		for _, v := range protoFiles {
			rel, err := filepath.Rel(gconfig.C.ProtoDir(), v)
			if err != nil {
				return nil, err
			}
			protoRelFiles = append(protoRelFiles, rel)
		}
		fds, err = protoParser.ParseFiles(protoRelFiles...)
		if err != nil {
			return nil, err
		}
	}

	rhis, err := jparser.Parse(g.config, fds, apiSpecs)
	if err != nil {
		return nil, err
	}

	var files []*GeneratedFile

	// 对 resource 进行排序
	resources := getResources(rhis)
	sort.Strings(resources)

	// 检查是否有 plugins
	plugins, err := plugin.GetPlugins()
	hasPlugins := err == nil && len(plugins) > 0

	// 分离未分组API
	ungroupedAPIs := getUngroupedAPIs(rhis)

	// gen clientset.go
	clientsetFile, err := g.genClientset(resources, hasPlugins, ungroupedAPIs)
	if err != nil {
		return nil, err
	}
	files = append(files, clientsetFile)

	// gen api types model
	apiTypesFile, err := g.genApiTypesModel(apiSpecs)
	if err != nil {
		return nil, err
	}
	files = append(files, apiTypesFile...)

	if len(protoFiles) > 0 {
		// gen pb model
		pbFiles, err := g.genPbTypesModel(protoFiles)
		if err != nil {
			return nil, err
		}
		files = append(files, pbFiles...)
	}

	for _, resource := range getResources(rhis) {
		resourcesFiles, err := g.genResources(rhis, resource)
		if err != nil {
			return nil, errors.Wrapf(err, "gen resources meet error. Resource: %s", resource)
		}
		files = append(files, resourcesFiles...)
	}

	// go mod file
	if g.config.GenModule {
		goModFile, err := g.genGoMod()
		if err != nil {
			return nil, err
		}
		files = append(files, goModFile)
	}

	// gen plugins
	if hasPlugins {
		pluginsFiles, err := g.genPlugins(plugins)
		if err != nil {
			return nil, err
		}
		files = append(files, pluginsFiles...)
	}

	return files, nil
}

func (g *Golang) genGoMod() (*GeneratedFile, error) {
	data, err := new.NewTemplateData()
	if err != nil {
		return nil, err
	}
	if gconfig.C.Gen.Sdk.GoVersion != "" {
		data["GoVersion"] = gconfig.C.Gen.Sdk.GoVersion
	}
	data["Module"] = gconfig.C.Gen.Sdk.GoModule

	template, err := templatex.ParseTemplate(filepath.Join("client", "client-go", "go.mod.tpl"), data, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "client-go", "go.mod.tpl"))))
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    "go.mod",
		Content: *bytes.NewBuffer(template),
	}, nil
}

func (g *Golang) genClientset(resources []string, hasPlugins bool, ungroupedAPIs []*vars.HTTPInterface) (*GeneratedFile, error) {
	clientGoBytes, err := templatex.ParseTemplate(filepath.Join("client", "client-go", "clientset.go.tpl"), map[string]any{
		"Package":       gconfig.C.Gen.Sdk.GoPackage,
		"Module":        gconfig.C.Gen.Sdk.GoModule,
		"Resources":     resources,
		"HasPlugins":    hasPlugins,
		"UngroupedAPIs": ungroupedAPIs,
		"GoImportPaths": g.genImports(ungroupedAPIs),
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "clientset.go.tpl")))
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    "clientset.go",
		Content: *bytes.NewBuffer(clientGoBytes),
	}, nil
}

func (g *Golang) genResources(rhis vars.ResourceHTTPInterfaceMap, resource string) ([]*GeneratedFile, error) {
	var resourceFiles []*GeneratedFile

	resourceGoBytes, err := templatex.ParseTemplate(filepath.Join("client", "client-go", "typed", "resource.go.tpl"), map[string]any{
		"GoModule":       gconfig.C.Gen.Sdk.GoModule,
		"Resource":       resource,
		"HTTPInterfaces": rhis[vars.Resource(resource)],
		"GoImportPaths":  g.genImports(rhis[vars.Resource(resource)]),
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "typed", "resource.go.tpl")))
	if err != nil {
		return nil, err
	}

	resourceFiles = append(resourceFiles, &GeneratedFile{
		Path:    filepath.Join("typed", strings.ToLower(resource), strings.ToLower(filepath.Base(resource))+".go"),
		Content: *bytes.NewBuffer(resourceGoBytes),
	})

	return resourceFiles, nil
}

func (g *Golang) genApiTypesModel(apiSpecs []*spec.ApiSpec) ([]*GeneratedFile, error) {
	var typesGoFiles []*GeneratedFile

	var allTypes []spec.Type

	// types 分组分文件夹生成
	for _, apiSpec := range apiSpecs {
		typesGoString, err := gogen.BuildTypes(apiSpec.Types)
		if err != nil {
			return nil, err
		}
		goPackage, ok := apiSpec.Info.Properties["go_package"]
		if ok {
			typesGoBytes, err := templatex.ParseTemplate("inner_types.go", map[string]any{
				"Types":   typesGoString,
				"Package": strings.ToLower(strings.ReplaceAll(goPackage, "/", "")),
			}, []byte(`// Code generated by jzero. DO NOT EDIT.
package {{.Package}}

import (
    "time"
)

var (
    _ = time.Now()
)

{{.Types}}`))
			if err != nil {
				return nil, err
			}

			logx.Debugf("get types Go Bytes: %v", string(typesGoBytes))
			source, err := gosimports.Process("", typesGoBytes, nil)
			if err != nil {
				return nil, err
			}
			typesGoFiles = append(typesGoFiles, &GeneratedFile{
				Path:    filepath.Join("model", goPackage, "types.go"),
				Content: *bytes.NewBuffer(source),
			})
		} else {
			allTypes = append(allTypes, apiSpec.Types...)
		}
	}

	if len(allTypes) > 0 {
		// 去除重复
		var realAllTypes []spec.Type
		exist := make(map[string]struct{})
		for _, v := range allTypes {
			if _, ok := exist[v.Name()]; ok {
				continue
			}
			realAllTypes = append(realAllTypes, v)
			exist[v.Name()] = struct{}{}
		}

		typesGoString, err := gogen.BuildTypes(realAllTypes)
		if err != nil {
			return nil, err
		}
		typesGoBytes, err := templatex.ParseTemplate("inner_types.go", map[string]any{
			"Types": typesGoString,
		}, []byte(`// Code generated by jzero. DO NOT EDIT.
package types

import (
    "time"
)

var (
    _ = time.Now()
)

{{.Types}}`))
		if err != nil {
			return nil, err
		}
		process, err := gosimports.Process("", typesGoBytes, nil)
		if err != nil {
			return nil, err
		}
		typesGoFiles = append(typesGoFiles, &GeneratedFile{
			Path:    filepath.Join("model", "types", "types.go"),
			Content: *bytes.NewBuffer(process),
		})
	}

	return typesGoFiles, nil
}

func (g *Golang) genPbTypesModel(protoFiles []string) ([]*GeneratedFile, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	for _, pf := range protoFiles {
		resp, err := execx.Run(fmt.Sprintf("protoc -I%s -I%s --go_out=%s %s", gconfig.C.ProtoDir(), filepath.Join(gconfig.C.ProtoDir(), "third_party"), tmpDir, pf), g.wd)
		if err != nil {
			return nil, errors.Errorf("err: [%v], resp: [%s]", err, resp)
		}
	}

	var generatedFiles []*GeneratedFile

	err = filepath.Walk(tmpDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}

		rel, err := filepath.Rel(tmpDir, filePath)
		if err != nil {
			return err
		}

		generatedFile := &GeneratedFile{
			Path:    filepath.Join("model", rel),
			Content: *bytes.NewBuffer(content),
		}

		generatedFiles = append(generatedFiles, generatedFile)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process directory: %v", err)
	}

	return generatedFiles, nil
}

func (g *Golang) genPlugins(plugins []plugin.Plugin) ([]*GeneratedFile, error) {
	var files []*GeneratedFile

	// 为每个 plugin 生成 model 和 typed 文件
	for _, p := range plugins {
		// 解析该 plugin 的 api 文件
		var apiSpecs []*spec.ApiSpec
		if pathx.FileExists(filepath.Join(p.Path, "desc", "api")) {
			apiFiles, err := jzerodesc.FindRouteApiFiles(filepath.Join(p.Path, "desc", "api"))
			if err != nil {
				continue
			}

			for _, apiFile := range apiFiles {
				apiSpec, err := apiparser.Parse(apiFile, nil)
				if err != nil {
					continue
				}
				apiSpecs = append(apiSpecs, apiSpec)
			}
		}

		if len(apiSpecs) > 0 {
			// 生成 plugin 的 model 文件
			pluginModelFiles, err := g.genPluginApiTypesModel(apiSpecs, p.Name)
			if err != nil {
				return nil, err
			}
			files = append(files, pluginModelFiles...)

			// 解析 plugin 的 api 生成 typed 文件
			rhis, err := jparser.Parse(g.config, nil, apiSpecs)
			if err != nil {
				return nil, err
			}

			// 生成 plugin 的 typed 文件
			for _, resource := range getResources(rhis) {
				pluginTypedFiles, err := g.genPluginResources(rhis, resource, p.Name)
				if err != nil {
					return nil, err
				}
				files = append(files, pluginTypedFiles...)
			}
		}
	}

	// 为每个插件生成单独的文件
	for _, p := range plugins {
		pluginGoFile, err := g.genSinglePluginGoFile(p)
		if err != nil {
			return nil, err
		}
		if pluginGoFile != nil {
			files = append(files, pluginGoFile)
		}
	}

	// 生成主 plugins.go 文件
	pluginsGoFile, err := g.genMainPluginsGoFile(plugins)
	if err != nil {
		return nil, err
	}
	files = append(files, pluginsGoFile)

	return files, nil
}

func (g *Golang) genPluginApiTypesModel(apiSpecs []*spec.ApiSpec, pluginName string) ([]*GeneratedFile, error) {
	var typesGoFiles []*GeneratedFile

	var allTypes []spec.Type

	// types 分组分文件夹生成
	for _, apiSpec := range apiSpecs {
		typesGoString, err := gogen.BuildTypes(apiSpec.Types)
		if err != nil {
			return nil, err
		}
		goPackage, ok := apiSpec.Info.Properties["go_package"]
		if ok {
			typesGoBytes, err := templatex.ParseTemplate("inner_types.go", map[string]any{
				"Types":   typesGoString,
				"Package": strings.ToLower(strings.ReplaceAll(goPackage, "/", "")),
			}, []byte(`// Code generated by jzero. DO NOT EDIT.
package {{.Package}}

import (
    "time"
)

var (
    _ = time.Now()
)

{{.Types}}`))
			if err != nil {
				return nil, err
			}

			source, err := gosimports.Process("", typesGoBytes, nil)
			if err != nil {
				return nil, err
			}
			typesGoFiles = append(typesGoFiles, &GeneratedFile{
				Path:    filepath.Join("plugins", pluginName, "model", goPackage, "types.go"),
				Content: *bytes.NewBuffer(source),
			})
		} else {
			allTypes = append(allTypes, apiSpec.Types...)
		}
	}

	if len(allTypes) > 0 {
		// 去除重复
		var realAllTypes []spec.Type
		exist := make(map[string]struct{})
		for _, v := range allTypes {
			if _, ok := exist[v.Name()]; ok {
				continue
			}
			realAllTypes = append(realAllTypes, v)
			exist[v.Name()] = struct{}{}
		}

		typesGoString, err := gogen.BuildTypes(realAllTypes)
		if err != nil {
			return nil, err
		}
		typesGoBytes, err := templatex.ParseTemplate("inner_types.go", map[string]any{
			"Types": typesGoString,
		}, []byte(`// Code generated by jzero. DO NOT EDIT.
package types

import (
    "time"
)

var (
    _ = time.Now()
)

{{.Types}}`))
		if err != nil {
			return nil, err
		}
		process, err := gosimports.Process("", typesGoBytes, nil)
		if err != nil {
			return nil, err
		}
		typesGoFiles = append(typesGoFiles, &GeneratedFile{
			Path:    filepath.Join("plugins", pluginName, "model", "types", "types.go"),
			Content: *bytes.NewBuffer(process),
		})
	}

	return typesGoFiles, nil
}

func (g *Golang) genPluginResources(rhis vars.ResourceHTTPInterfaceMap, resource string, pluginName string) ([]*GeneratedFile, error) {
	var resourceFiles []*GeneratedFile

	resourceGoBytes, err := templatex.ParseTemplate(filepath.Join("client", "client-go", "typed", "resource.go.tpl"), map[string]any{
		"GoModule":       gconfig.C.Gen.Sdk.GoModule,
		"Resource":       resource,
		"HTTPInterfaces": rhis[vars.Resource(resource)],
		"GoImportPaths":  g.genPluginImports(rhis[vars.Resource(resource)], pluginName),
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "typed", "resource.go.tpl")))
	if err != nil {
		return nil, err
	}

	resourceFiles = append(resourceFiles, &GeneratedFile{
		Path:    filepath.Join("plugins", pluginName, "typed", strings.ToLower(resource), strings.ToLower(filepath.Base(resource))+".go"),
		Content: *bytes.NewBuffer(resourceGoBytes),
	})

	return resourceFiles, nil
}

func (g *Golang) genPluginImports(infs []*vars.HTTPInterface, pluginName string) []string {
	var imports []string
	for _, inf := range infs {
		if inf.Request != nil && inf.Request.Package != "" {
			imports = append(imports, fmt.Sprintf("%s/plugins/%s/model/%s", gconfig.C.Gen.Sdk.GoModule, pluginName, inf.Request.Package))
		}
		if inf.Response != nil && inf.Response.Package != "" {
			imports = append(imports, fmt.Sprintf("%s/plugins/%s/model/%s", gconfig.C.Gen.Sdk.GoModule, pluginName, inf.Response.Package))
		}
	}
	return imports
}

// genSinglePluginGoFile 为单个插件生成文件
func (g *Golang) genSinglePluginGoFile(p plugin.Plugin) (*GeneratedFile, error) {
	var apiSpecs []*spec.ApiSpec
	if pathx.FileExists(filepath.Join(p.Path, "desc", "api")) {
		apiFiles, err := jzerodesc.FindRouteApiFiles(filepath.Join(p.Path, "desc", "api"))
		if err != nil {
			return nil, nil
		}

		for _, apiFile := range apiFiles {
			apiSpec, err := apiparser.Parse(apiFile, nil)
			if err != nil {
				continue
			}
			apiSpecs = append(apiSpecs, apiSpec)
		}
	}

	if len(apiSpecs) == 0 {
		return nil, nil
	}

	rhis, err := jparser.Parse(g.config, nil, apiSpecs)
	if err != nil {
		return nil, nil
	}

	resources := getResources(rhis)
	sort.Strings(resources)
	ungroupedAPIs := getUngroupedAPIs(rhis)

	// 使用单个插件模板
	pluginGoBytes, err := templatex.ParseTemplate(filepath.Join("client", "client-go", "plugin.go.tpl"), map[string]any{
		"Module":        gconfig.C.Gen.Sdk.GoModule,
		"PluginName":    p.Name,
		"Resources":     resources,
		"UngroupedAPIs": ungroupedAPIs,
		"GoImportPaths": g.genPluginImports(ungroupedAPIs, p.Name),
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "plugin.go.tpl")))
	if err != nil {
		return nil, err
	}

	source, err := gosimports.Process("", pluginGoBytes, nil)
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    filepath.Join("plugins", p.Name+".go"),
		Content: *bytes.NewBuffer(source),
	}, nil
}

func (g *Golang) genMainPluginsGoFile(plugins []plugin.Plugin) (*GeneratedFile, error) {
	// 只收集插件名称，具体实现在各自的文件中
	var pluginNames []string
	for _, p := range plugins {
		// 检查插件是否有API文件
		if pathx.FileExists(filepath.Join(p.Path, "desc", "api")) {
			apiFiles, err := jzerodesc.FindRouteApiFiles(filepath.Join(p.Path, "desc", "api"))
			if err == nil && len(apiFiles) > 0 {
				pluginNames = append(pluginNames, p.Name)
			}
		}
	}

	if len(pluginNames) == 0 {
		// 如果没有插件，生成一个空的 plugins.go
		pluginNames = []string{}
	}

	// 使用主插件模板文件生成 plugins.go
	pluginsGoBytes, err := templatex.ParseTemplate(filepath.Join("client", "client-go", "plugins.go.tpl"), map[string]any{
		"Module":      gconfig.C.Gen.Sdk.GoModule,
		"PluginNames": pluginNames,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "plugins.go.tpl")))
	if err != nil {
		return nil, err
	}

	source, err := gosimports.Process("", pluginsGoBytes, nil)
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    filepath.Join("plugins", "plugins.go"),
		Content: *bytes.NewBuffer(source),
	}, nil
}

func (g *Golang) genImports(infs []*vars.HTTPInterface) []string {
	var imports []string
	for _, inf := range infs {
		if inf.Request != nil && inf.Request.Package != "" {
			imports = append(imports, fmt.Sprintf("%s/model/%s", gconfig.C.Gen.Sdk.GoModule, inf.Request.Package))
		}
		if inf.Response != nil && inf.Response.Package != "" {
			imports = append(imports, fmt.Sprintf("%s/model/%s", gconfig.C.Gen.Sdk.GoModule, inf.Response.Package))
		}
	}
	return imports
}
