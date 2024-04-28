package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/jaronnie/jzero/cmd/gen"
	"github.com/jaronnie/jzero/cmd/gensdk/config"
	"github.com/jaronnie/jzero/cmd/gensdk/jparser"
	"github.com/jaronnie/jzero/cmd/gensdk/vars"
	"github.com/jaronnie/jzero/daemon/pkg/templatex"
	"github.com/jaronnie/jzero/embeded"
	"github.com/jhump/protoreflect/desc/protoparse"
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

	// parse proto
	var protoParser protoparse.Parser
	protoParser.ImportPaths = []string{filepath.Join("daemon", "desc", "proto")}

	// parse api
	var apiSpecs []*spec.ApiSpec
	apiSpec, err := apiparser.Parse(filepath.Join("daemon", "desc", "api", g.config.APP+".api"))
	if err != nil {
		return nil, err
	}
	apiSpecs = append(apiSpecs, apiSpec)

	protoFiles, err := gen.GetProtoFilenames(wd)
	if err != nil {
		return nil, err
	}
	fds, err := protoParser.ParseFiles(protoFiles...)
	if err != nil {
		return nil, err
	}

	rhis, err := jparser.Parse(g.config, fds, apiSpecs)
	if err != nil {
		return nil, err
	}

	var files []*GeneratedFile

	// gen clientset.go
	clientGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":    g.config.APP,
		"Module": g.config.Module,
		"Scopes": getScopes(rhis),
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "clientset.go.tpl")))
	if err != nil {
		return nil, err
	}
	files = append(files, &GeneratedFile{
		Path:    "clientset.go",
		Content: *bytes.NewBuffer(clientGoBytes),
	})

	// gen rest frame
	files = append(files, &GeneratedFile{
		Path:    filepath.Join("rest", "client.go"),
		Content: *bytes.NewBuffer(embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "rest", "client.go.tpl"))),
	})

	files = append(files, &GeneratedFile{
		Path:    filepath.Join("rest", "option.go"),
		Content: *bytes.NewBuffer(embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "rest", "option.go.tpl"))),
	})

	files = append(files, &GeneratedFile{
		Path:    filepath.Join("rest", "request.go"),
		Content: *bytes.NewBuffer(embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "rest", "request.go.tpl"))),
	})

	// gen typed/direct_client.go
	directClientGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Module": g.config.Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "direct_client.go.tpl")))
	if err != nil {
		return nil, err
	}
	files = append(files, &GeneratedFile{
		Path:    filepath.Join("typed", "direct_client.go"),
		Content: *bytes.NewBuffer(directClientGoBytes),
	})

	for _, scope := range getScopes(rhis) {
		// gen typed/scope_client.go
		scopeClientGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
			"Scope":     scope,
			"Module":    g.config.Module,
			"Resources": getScopeResources(rhis[vars.Scope(scope)]),
		}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "scope_client.go.tpl")))
		if err != nil {
			return nil, err
		}
		files = append(files, &GeneratedFile{
			Path:    filepath.Join("typed", scope, scope+"_client.go"),
			Content: *bytes.NewBuffer(scopeClientGoBytes),
		})

		// gen api types model
		apiTypes, err := g.genApiTypesModel(apiSpec.Types)
		apiTypesGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
			"Types": apiTypes,
		}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "model", "types", "scope_types.go.tpl")))
		if err != nil {
			return nil, err
		}
		files = append(files, &GeneratedFile{
			Path:    filepath.Join("model", scope, "types", "types.go"),
			Content: *bytes.NewBuffer(apiTypesGoBytes),
		})

		// gen pb model
		pbFiles, err := g.genPbTypesModel()
		if err != nil {
			return nil, err
		}
		files = append(files, pbFiles...)

		for _, resource := range getScopeResources(rhis[vars.Scope(scope)]) {
			// gen typed/resource_expansion.go
			resourceExpansionGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
				"Module":   g.config.Module,
				"Scope":    scope,
				"Resource": resource,
			}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "resource_expansion.go.tpl")))
			if err != nil {
				return nil, err
			}
			files = append(files, &GeneratedFile{
				Path:    filepath.Join("typed", scope, resource+"_expansion.go"),
				Content: *bytes.NewBuffer(resourceExpansionGoBytes),
			})

			resourceExpansionGoBytes, err = templatex.ParseTemplate(map[string]interface{}{
				"Module":   g.config.Module,
				"Scope":    scope,
				"Resource": resource,
			}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "resource_expansion.go.tpl")))
			if err != nil {
				return nil, err
			}
			files = append(files, &GeneratedFile{
				Path:    filepath.Join("typed", scope, resource+"_expansion.go"),
				Content: *bytes.NewBuffer(resourceExpansionGoBytes),
			})

			// gen typed/resource.go
			resourceGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
				"GoModule":           g.config.Module,
				"Scope":              scope,
				"Resource":           resource,
				"HTTPInterfaces":     rhis[vars.Scope(scope)][vars.Resource(resource)],
				"IsWarpHTTPResponse": true,
				"GoImportPaths":      g.genImports(rhis[vars.Scope(scope)][vars.Resource(resource)]),
			}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "resource.go.tpl")))
			if err != nil {
				return nil, err
			}
			files = append(files, &GeneratedFile{
				Path:    filepath.Join("typed", scope, resource+".go"),
				Content: *bytes.NewBuffer(resourceGoBytes),
			})
		}
	}

	// go mod file
	goModFile, err := g.genGoMod()
	if err != nil {
		return nil, err
	}
	files = append(files, goModFile)

	return files, nil
}

func (g *Golang) genGoMod() (*GeneratedFile, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	resp, err := execx.Run(fmt.Sprintf("go mod init %s", g.config.Module), tmpDir)
	if err != nil {
		return nil, errors.Errorf("err: [%v], resp: [%s]", err, resp)
	}

	goModBytes, err := os.ReadFile(filepath.Join(tmpDir, "go.mod"))
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Skip:    true,
		Path:    "go.mod",
		Content: *bytes.NewBuffer(goModBytes),
	}, nil
}

func (g *Golang) genApiTypesModel(types []spec.Type) (string, error) {
	return gogen.BuildTypes(types)
}

func (g *Golang) genPbTypesModel() ([]*GeneratedFile, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	resp, err := execx.Run(fmt.Sprintf("protoc -I./daemon/desc/proto --go_out=%s daemon/desc/proto/*.proto", tmpDir), g.wd)
	if err != nil {
		return nil, errors.Errorf("err: [%v], resp: [%s]", err, resp)
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
			Path:    filepath.Join("model", g.config.APP, rel),
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
		imports = append(imports, fmt.Sprintf("%s/model/%s/%s", g.config.Module, g.config.APP, inf.RequestBody.Package))
		imports = append(imports, fmt.Sprintf("%s/model/%s/%s", g.config.Module, g.config.APP, inf.ResponseBody.Package))
	}
	return imports
}
