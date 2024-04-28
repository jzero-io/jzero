package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jaronnie/jzero/cmd/gensdk/vars"

	"github.com/jaronnie/jzero/cmd/gen"
	"github.com/jaronnie/jzero/cmd/gensdk/config"
	"github.com/jaronnie/jzero/daemon/pkg/templatex"
	"github.com/jaronnie/jzero/embeded"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"

	"github.com/jaronnie/jzero/cmd/gensdk/jparser"
	"github.com/jhump/protoreflect/desc/protoparse"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
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
}

func (g *Golang) Gen() ([]*GeneratedFile, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

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

	var apiGoImportPaths []string

	protoFiles, err := gen.GetProtoFilenames(wd)
	if err != nil {
		return nil, err
	}
	fds, err := protoParser.ParseFiles(protoFiles...)
	if err != nil {
		return nil, err
	}

	apiGoImportPaths = append(apiGoImportPaths, fmt.Sprintf("%s/model/types", g.config.Module))

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

		for _, resource := range getScopeResources(rhis[vars.Scope(scope)]) {
			// TODO get go imports
			var protoGoImportPaths []string
			protoGoImportPaths = append(protoGoImportPaths, fmt.Sprintf("%s/model/%s", g.config.Module, rhis[vars.Scope(scope)][vars.Resource(resource)][0].RequestBody.Package))

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
				"GoImportPaths":      protoGoImportPaths,
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

	// go mod init
	_, err = execx.Run(fmt.Sprintf("go mod init %s", g.config.Module), g.config.Dir)

	return files, nil
}
