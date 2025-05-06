package generator

import (
	"bytes"
	"fmt"
	goformat "go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	gconfig "github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/jparser"
	"github.com/jzero-io/jzero/internal/gen/gensdk/vars"
	"github.com/jzero-io/jzero/internal/new"
	jzerodesc "github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/osx"
	"github.com/jzero-io/jzero/pkg/templatex"
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
			apiSpec, err := apiparser.Parse(v)
			if err != nil {
				return nil, err
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

	// gen clientset.go
	clientsetFiles, err := g.genClientSets(getScopes(rhis))
	if err != nil {
		return nil, err
	}
	files = append(files, clientsetFiles...)

	// gen direct_client
	directClientFiles, err := g.genDirectClients()
	if err != nil {
		return nil, err
	}
	files = append(files, directClientFiles...)

	for _, scope := range getScopes(rhis) {
		// 对 resource 进行排序
		resources := getScopeResources(rhis[vars.Scope(scope)])
		sort.Strings(resources)
		scopeClientFiles, err := g.genScopeClients(scope, resources)
		if err != nil {
			return nil, err
		}
		files = append(files, scopeClientFiles...)

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

		for _, resource := range getScopeResources(rhis[vars.Scope(scope)]) {
			scopeResourcesFiles, err := g.genScopeResources(rhis, scope, resource)
			if err != nil {
				return nil, errors.Wrapf(err, "gen scope resources meet error. Resource: %s", resource)
			}
			files = append(files, scopeResourcesFiles...)
		}
	}

	// go mod file
	if g.config.GenModule {
		goModFile, err := g.genGoMod()
		if err != nil {
			return nil, err
		}
		files = append(files, goModFile)
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

	template, err := templatex.ParseTemplate(data, embeded.ReadTemplateFile(filepath.ToSlash(filepath.Join("client", "client-go", "go.mod.tpl"))))
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    "go.mod",
		Content: *bytes.NewBuffer(template),
	}, nil
}

func (g *Golang) genClientSets(scopes []string) ([]*GeneratedFile, error) {
	var clientSetFiles []*GeneratedFile

	clientGoBytes, err := templatex.ParseTemplate(map[string]any{
		"Package": gconfig.C.Gen.Sdk.GoPackage,
		"Module":  gconfig.C.Gen.Sdk.GoModule,
		"Scopes":  scopes,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "clientset.go.tpl")))
	if err != nil {
		return nil, err
	}
	clientSetFiles = append(clientSetFiles, &GeneratedFile{
		Path:    "clientset.go",
		Content: *bytes.NewBuffer(clientGoBytes),
	})

	return clientSetFiles, nil
}

func (g *Golang) genDirectClients() ([]*GeneratedFile, error) {
	var directClientFiles []*GeneratedFile

	directClientGoBytes, err := templatex.ParseTemplate(map[string]any{
		"Module": gconfig.C.Gen.Sdk.GoModule,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "typed", "direct_client.go.tpl")))
	if err != nil {
		return nil, err
	}
	directClientFiles = append(directClientFiles, &GeneratedFile{
		Path:    filepath.Join("typed", "direct_client.go"),
		Content: *bytes.NewBuffer(directClientGoBytes),
	})

	return directClientFiles, nil
}

func (g *Golang) genScopeClients(scope string, resources []string) ([]*GeneratedFile, error) {
	var scopeClientFiles []*GeneratedFile

	scopeClientGoBytes, err := templatex.ParseTemplate(map[string]any{
		"Scope":     scope,
		"Module":    gconfig.C.Gen.Sdk.GoModule,
		"Resources": resources,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "typed", "scope_client.go.tpl")))
	if err != nil {
		return nil, err
	}

	scopeClientFiles = append(scopeClientFiles, &GeneratedFile{
		Path:    filepath.Join("typed", strings.ToLower(scope), strings.ToLower(scope)+"_client.go"),
		Content: *bytes.NewBuffer(scopeClientGoBytes),
	})

	return scopeClientFiles, nil
}

func (g *Golang) genScopeResources(rhis vars.ScopeResourceHTTPInterfaceMap, scope, resource string) ([]*GeneratedFile, error) {
	var scopeResourceFiles []*GeneratedFile

	// resource_expansion.go
	resourceExpansionGoBytes, err := templatex.ParseTemplate(map[string]any{
		"Module":   gconfig.C.Gen.Sdk.GoModule,
		"Scope":    scope,
		"Resource": resource,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "typed", "resource_expansion.go.tpl")))
	if err != nil {
		return nil, err
	}
	scopeResourceFiles = append(scopeResourceFiles, &GeneratedFile{
		Path:    filepath.Join("typed", strings.ToLower(scope), strings.ToLower(resource)+"_expansion.go"),
		Content: *bytes.NewBuffer(resourceExpansionGoBytes),
	})

	resourceGoBytes, err := templatex.ParseTemplate(map[string]any{
		"GoModule":           gconfig.C.Gen.Sdk.GoModule,
		"Scope":              scope,
		"Resource":           resource,
		"HTTPInterfaces":     rhis[vars.Scope(scope)][vars.Resource(resource)],
		"IsWrapHTTPResponse": gconfig.C.Gen.Sdk.WrapResponse,
		"GoImportPaths":      g.genImports(rhis[vars.Scope(scope)][vars.Resource(resource)]),
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "typed", "resource.go.tpl")))
	if err != nil {
		return nil, err
	}

	// resourceGoFormatBytes, err := gosimports.Process("", resourceGoBytes, &gosimports.Options{Comments: true})
	// if err != nil {
	//	return nil, errors.Errorf("format resource.go meet error: %s", err)
	// }

	scopeResourceFiles = append(scopeResourceFiles, &GeneratedFile{
		Path:    filepath.Join("typed", strings.ToLower(scope), strings.ToLower(resource)+".go"),
		Content: *bytes.NewBuffer(resourceGoBytes),
	})

	return scopeResourceFiles, nil
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
			typesGoBytes, err := templatex.ParseTemplate(map[string]any{
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

			source, err := goformat.Source(typesGoBytes)
			if err != nil {
				return nil, err
			}
			typesGoFiles = append(typesGoFiles, &GeneratedFile{
				Path:    filepath.Join("model", strings.ToLower(gconfig.C.Gen.Sdk.Scope), goPackage, "types.go"),
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
		typesGoBytes, err := templatex.ParseTemplate(map[string]any{
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
		source, err := goformat.Source(typesGoBytes)
		if err != nil {
			return nil, err
		}
		typesGoFiles = append(typesGoFiles, &GeneratedFile{
			Path:    filepath.Join("model", strings.ToLower(gconfig.C.Gen.Sdk.Scope), "types", "types.go"),
			Content: *bytes.NewBuffer(source),
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
			Path:    filepath.Join("model", strings.ToLower(gconfig.C.Gen.Sdk.Scope), rel),
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

func (g *Golang) genImports(infs []*vars.HTTPInterface) []string {
	var imports []string
	for _, inf := range infs {
		if inf.Request != nil && inf.Request.Package != "" {
			imports = append(imports, fmt.Sprintf("%s/model/%s/%s", gconfig.C.Gen.Sdk.GoModule, strings.ToLower(gconfig.C.Gen.Sdk.Scope), inf.Request.Package))
		}
		if inf.Response != nil && inf.Response.Package != "" {
			imports = append(imports, fmt.Sprintf("%s/model/%s/%s", gconfig.C.Gen.Sdk.GoModule, strings.ToLower(gconfig.C.Gen.Sdk.Scope), inf.Response.Package))
		}
	}
	return imports
}
