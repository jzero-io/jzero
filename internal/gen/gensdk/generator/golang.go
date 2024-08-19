package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jzero-io/jzero/internal/new"
	"github.com/rinchsan/gosimports"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/jparser"
	"github.com/jzero-io/jzero/internal/gen/gensdk/vars"

	"github.com/jhump/protoreflect/desc"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/pkg/errors"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/pkg/templatex"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
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

	if pathx.FileExists(g.config.ApiDir) {
		mainApiFilePath, isDelete, err := gen.GetMainApiFilePath(g.config.ApiDir)
		if err != nil {
			return nil, err
		}

		apiSpec, err := apiparser.Parse(mainApiFilePath)
		if err != nil {
			return nil, err
		}

		if isDelete {
			_ = os.Remove(mainApiFilePath)
		}

		apiSpecs = append(apiSpecs, apiSpec)
	}

	var protoFiles []string

	if pathx.FileExists(g.config.ProtoDir) {
		protoFiles, err = gen.GetProtoFilepath(g.config.ProtoDir)
		if err != nil {
			return nil, err
		}

	}

	var fds []*desc.FileDescriptor

	// parse proto
	var protoParser protoparse.Parser
	if len(protoFiles) > 0 {
		protoParser.ImportPaths = []string{g.config.ProtoDir}
		var protoRelFiles []string
		for _, v := range protoFiles {
			rel, err := filepath.Rel(g.config.ProtoDir, v)
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
		scopeClientFiles, err := g.genScopeClients(scope, getScopeResources(rhis[vars.Scope(scope)]))
		if err != nil {
			return nil, err
		}
		files = append(files, scopeClientFiles...)

		// gen api types model
		if len(apiSpecs) > 0 {
			apiTypesFile, err := g.genApiTypesModel(apiSpecs[0].Types)
			if err != nil {
				return nil, err
			}
			files = append(files, apiTypesFile)
		}

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
				return nil, err
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
	if g.config.GoVersion != "" {
		data["GoVersion"] = g.config.GoVersion
	}
	data["Module"] = g.config.GoModule

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

	clientGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Package": g.config.GoPackage,
		"Module":  g.config.GoModule,
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

	directClientGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": g.config.GoModule,
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

	scopeClientGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Scope":     scope,
		"Module":    g.config.GoModule,
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

func (g *Golang) genScopeResources(rhis vars.ScopeResourceHTTPInterfaceMap, scope string, resource string) ([]*GeneratedFile, error) {
	var scopeResourceFiles []*GeneratedFile

	// resource_expansion.go
	resourceExpansionGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Module":   g.config.GoModule,
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

	resourceGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"GoModule":           g.config.GoModule,
		"Scope":              scope,
		"Resource":           resource,
		"HTTPInterfaces":     rhis[vars.Scope(scope)][vars.Resource(resource)],
		"IsWrapHTTPResponse": g.config.WrapResponse,
		"GoImportPaths":      g.genImports(rhis[vars.Scope(scope)][vars.Resource(resource)]),
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "typed", "resource.go.tpl")))
	if err != nil {
		return nil, err
	}

	resourceGoFormatBytes, err := gosimports.Process("", resourceGoBytes, &gosimports.Options{Comments: true})
	if err != nil {
		return nil, err
	}

	scopeResourceFiles = append(scopeResourceFiles, &GeneratedFile{
		Path:    filepath.Join("typed", strings.ToLower(scope), strings.ToLower(resource)+".go"),
		Content: *bytes.NewBuffer(resourceGoFormatBytes),
	})

	return scopeResourceFiles, nil
}

func (g *Golang) genApiTypesModel(types []spec.Type) (*GeneratedFile, error) {
	typesGoString, err := gogen.BuildTypes(types)
	if err != nil {
		return nil, err
	}

	typesGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Types": typesGoString,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-go", "model", "types", "scope_types.go.tpl")))
	if err != nil {
		return nil, err
	}
	return &GeneratedFile{
		Path:    filepath.Join("model", strings.ToLower(g.config.Scope), "types", "types.go"),
		Content: *bytes.NewBuffer(typesGoBytes),
	}, nil
}

func (g *Golang) genPbTypesModel(protoFiles []string) ([]*GeneratedFile, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	for _, pf := range protoFiles {
		resp, err := execx.Run(fmt.Sprintf("protoc -I%s --go_out=%s %s", g.config.ProtoDir, tmpDir, pf), g.wd)
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
			Path:    filepath.Join("model", strings.ToLower(g.config.Scope), rel),
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
			imports = append(imports, fmt.Sprintf("%s/model/%s/%s", g.config.GoModule, strings.ToLower(g.config.Scope), inf.Request.Package))
		}
		if inf.Response != nil && inf.Response.Package != "" {
			imports = append(imports, fmt.Sprintf("%s/model/%s/%s", g.config.GoModule, strings.ToLower(g.config.Scope), inf.Response.Package))
		}
	}
	return imports
}
