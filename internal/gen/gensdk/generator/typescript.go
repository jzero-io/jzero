package generator

import (
	"bytes"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/tsgen"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	gconfig "github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/embeded"
	"github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/jparser"
	"github.com/jzero-io/jzero/internal/gen/gensdk/vars"
	jzerodesc "github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/templatex"
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
}

func (t *Typescript) Gen() ([]*GeneratedFile, error) {
	// parse api
	var apiSpecs []*spec.ApiSpec

	if pathx.FileExists(gconfig.C.ApiDir()) {
		files, err := jzerodesc.FindApiFiles(gconfig.C.ApiDir())
		if err != nil {
			return nil, err
		}
		for _, v := range files {
			apiSpec, err := apiparser.Parse(v)
			if err != nil {
				return nil, err
			}
			apiSpecs = append(apiSpecs, apiSpec)
		}
	}

	protoFiles, err := jzerodesc.GetProtoFilepath(gconfig.C.ProtoDir())
	if err != nil {
		return nil, err
	}

	var fds []*desc.FileDescriptor

	// parse proto
	var protoParser protoparse.Parser
	if len(protoFiles) > 0 {
		protoParser.ImportPaths = []string{gconfig.C.ProtoDir()}
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
	clientBytes, err := templatex.ParseTemplate(map[string]any{
		"APP":    gconfig.C.Gen.Sdk.Scope,
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
	packageJsonBytes, err := templatex.ParseTemplate(map[string]any{
		"APP": gconfig.C.Gen.Sdk.Scope,
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
	scopeClientBytes, err := templatex.ParseTemplate(map[string]any{
		"Scope":     scope,
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
	requestBytes, err := templatex.ParseTemplate(map[string]any{}, embeded.ReadTemplateFile(filepath.Join("client", "client-ts", "rest", "request.ts.tpl")))
	if err != nil {
		return nil, err
	}

	return &GeneratedFile{
		Path:    filepath.Join("rest", "request.ts"),
		Content: *bytes.NewBuffer(requestBytes),
	}, nil
}

func (t *Typescript) genScopeResources(rhis vars.ScopeResourceHTTPInterfaceMap, scope, resource string) ([]*GeneratedFile, error) {
	var scopeResourceFiles []*GeneratedFile

	resourceGoBytes, err := templatex.ParseTemplate(map[string]any{
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
		Path:    filepath.Join("model", gconfig.C.Gen.Sdk.Scope, "types", "types.ts"),
		Content: *bytes.NewBuffer([]byte(typesGoString)),
	}, nil
}
