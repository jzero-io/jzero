package generator

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/generator/tsgen"
	"github.com/jzero-io/jzero/internal/gen/gensdk/jparser"
	"github.com/jzero-io/jzero/internal/gen/gensdk/vars"
	"github.com/jzero-io/jzero/pkg/templatex"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func init() {
	Register("ts", func(config config.Config) (Generator, error) {
		return &Typescript{
			config: &config,
		}, nil
	})
}

type Typescript struct {
	config *config.Config

	wd string
}

func (t *Typescript) Gen() ([]*GeneratedFile, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	t.wd = wd

	// parse api
	var apiSpecs []*spec.ApiSpec

	if pathx.FileExists(t.config.ApiDir) {
		mainApiFilePath, isDelete, err := gen.GetMainApiFilePath(t.config.ApiDir)
		if isDelete {
			defer os.Remove(mainApiFilePath)
		}
		if err != nil {
			return nil, err
		}
		apiSpec, err := apiparser.Parse(mainApiFilePath)
		if err != nil {
			return nil, err
		}
		if mainApiFilePath != filepath.Join(t.config.ApiDir, "main.api") {
			os.Remove(mainApiFilePath)
		}

		apiSpecs = append(apiSpecs, apiSpec)
	}

	protoFiles, err := gen.GetProtoFilepath(t.config.ProtoDir)
	if err != nil {
		return nil, err
	}

	var fds []*desc.FileDescriptor

	// parse proto
	var protoParser protoparse.Parser
	if len(protoFiles) > 0 {
		protoParser.ImportPaths = []string{t.config.ProtoDir}
		var protoRelFiles []string
		for _, v := range protoFiles {
			rel, err := filepath.Rel(t.config.ProtoDir, v)
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

	rhis, err := jparser.Parse(t.config, fds, apiSpecs)
	if err != nil {
		return nil, err
	}

	var files []*GeneratedFile

	clientSetFile, err := t.genClientSet(getScopes(rhis))
	if err != nil {
		return nil, err
	}
	files = append(files, clientSetFile)

	packageJsonFile, err := t.genPackageJson(getScopes(rhis))
	if err != nil {
		return nil, err
	}
	files = append(files, packageJsonFile)

	restFile, err := t.genRest()
	if err != nil {
		return nil, err
	}
	files = append(files, restFile)

	for _, scope := range getScopes(rhis) {
		scopeClientFile, err := t.genScopeClient(scope, getScopeResources(rhis[vars.Scope(scope)]))
		if err != nil {
			return nil, err
		}
		files = append(files, scopeClientFile)

		// gen api types model
		if len(apiSpecs) > 0 {
			apiTypesFile, err := t.genApiTypesModel(apiSpecs[0].Types)
			if err != nil {
				return nil, err
			}
			files = append(files, apiTypesFile)
		}

		for _, resource := range getScopeResources(rhis[vars.Scope(scope)]) {
			scopeResourcesFiles, err := t.genScopeResources(rhis, scope, resource)
			if err != nil {
				return nil, err
			}
			files = append(files, scopeResourcesFiles...)
		}
	}

	return files, nil
}

func (t *Typescript) genClientSet(scopes []string) (*GeneratedFile, error) {
	clientBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":    t.config.APP,
		"Scopes": scopes,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-ts", "index.ts.tpl")))
	if err != nil {
		return nil, err
	}
	return &GeneratedFile{
		Path:    "index.ts",
		Content: *bytes.NewBuffer(clientBytes),
	}, nil
}

func (t *Typescript) genPackageJson(scopes []string) (*GeneratedFile, error) {
	packageJsonBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"APP": t.config.APP,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-ts", "package.json.tpl")))
	if err != nil {
		return nil, err
	}
	return &GeneratedFile{
		Path:    "package.json",
		Content: *bytes.NewBuffer(packageJsonBytes),
	}, nil
}

func (t *Typescript) genScopeClient(scope string, resources []string) (*GeneratedFile, error) {
	scopeClientBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Scope":     scope,
		"Module":    t.config.Module,
		"Resources": resources,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-ts", "typed", "scope_client.ts.tpl")))
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    filepath.Join("typed", scope, scope+"_client.ts"),
		Content: *bytes.NewBuffer(scopeClientBytes),
	}, nil
}

func (t *Typescript) genRest() (*GeneratedFile, error) {
	requestBytes, err := templatex.ParseTemplate(map[string]interface{}{}, embeded.ReadTemplateFile(filepath.Join("client", "client-ts", "rest", "request.ts.tpl")))
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    filepath.Join("rest", "request.ts"),
		Content: *bytes.NewBuffer(requestBytes),
	}, nil
}

func (t *Typescript) genScopeResources(rhis vars.ScopeResourceHTTPInterfaceMap, scope string, resource string) ([]*GeneratedFile, error) {
	var scopeResourceFiles []*GeneratedFile

	resourceGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"HTTPInterfaces": rhis[vars.Scope(scope)][vars.Resource(resource)],
		"Resource":       resource,
	}, embeded.ReadTemplateFile(filepath.Join("client", "client-ts", "typed", "resource.ts.tpl")))
	if err != nil {
		return nil, err
	}
	scopeResourceFiles = append(scopeResourceFiles, &GeneratedFile{
		Path:    filepath.Join("typed", scope, resource+".ts"),
		Content: *bytes.NewBuffer(resourceGoBytes),
	})

	return scopeResourceFiles, nil
}

func (t *Typescript) genApiTypesModel(types []spec.Type) (*GeneratedFile, error) {
	typesGoString, err := tsgen.BuildTypes(types)
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    filepath.Join("model", t.config.APP, "types", "types.ts"),
		Content: *bytes.NewBuffer([]byte(typesGoString)),
	}, nil
}
