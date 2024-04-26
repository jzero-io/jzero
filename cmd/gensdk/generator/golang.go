package generator

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/jaronnie/jzero/daemon/pkg/templatex"
	"github.com/jaronnie/jzero/embeded"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"

	"github.com/jaronnie/jzero/cmd/gensdk/jparser"
	"github.com/jhump/protoreflect/desc/protoparse"
	apiparser "github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

func init() {
	Register("go", func(target Target) (Generator, error) {
		return &Golang{
			target: &target,
		}, nil
	})
}

type Golang struct {
	target *Target
}

func (g *Golang) Gen() ([]*GeneratedFile, error) {
	// parse proto
	var protoParser protoparse.Parser
	protoParser.ImportPaths = []string{filepath.Join("daemon", "desc", "proto")}

	// parse api
	var apiSpecs []*spec.ApiSpec
	apiSpec, err := apiparser.Parse(filepath.Join("daemon", "desc", "api", g.target.APP+".api"))
	if err != nil {
		return nil, err
	}
	apiSpecs = append(apiSpecs, apiSpec)

	fds, err := protoParser.ParseFiles("credential.proto", "machine.proto")
	if err != nil {
		return nil, err
	}

	rhis, err := jparser.Parse(fds, apiSpecs)
	if err != nil {
		return nil, err
	}
	fmt.Println(rhis)

	var files []*GeneratedFile

	// gen clientset.go
	clientGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"APP":    g.target.APP,
		"Module": g.target.Module,
		"Scopes": []string{"jzero"},
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "clientset.go.tpl")))
	if err != nil {
		return nil, err
	}

	// gen rest frame
	files = append(files, &GeneratedFile{
		Path:    "clientset.go",
		Content: *bytes.NewBuffer(clientGoBytes),
	})

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
		"Module": g.target.Module,
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "direct_client.go.tpl")))
	if err != nil {
		return nil, err
	}
	files = append(files, &GeneratedFile{
		Path:    filepath.Join("typed", "direct_client.go"),
		Content: *bytes.NewBuffer(directClientGoBytes),
	})

	// gen typed/scope_client.go
	scopeClientGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Scope":     "jzero",
		"Module":    g.target.Module,
		"Resources": []string{"credential"},
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "scope_client.go.tpl")))
	if err != nil {
		return nil, err
	}
	files = append(files, &GeneratedFile{
		Path:    filepath.Join("typed", "jzero", "jzero_client.go"),
		Content: *bytes.NewBuffer(scopeClientGoBytes),
	})

	// gen typed/resource_expansion.go
	resourceExpansionGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"Module":   g.target.Module,
		"Scope":    "jzero",
		"Resource": "credential",
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "resource_expansion.go.tpl")))
	if err != nil {
		return nil, err
	}
	files = append(files, &GeneratedFile{
		Path:    filepath.Join("typed", "jzero", "credential_expansion.go"),
		Content: *bytes.NewBuffer(resourceExpansionGoBytes),
	})

	// gen typed/resource.go
	resourceGoBytes, err := templatex.ParseTemplate(map[string]interface{}{
		"GoModule": g.target.Module,
		"Scope":    "jzero",
		"Resource": "credential",
	}, embeded.ReadTemplateFile(filepath.Join("jzero", "client", "client-go", "typed", "resource.go.tpl")))
	if err != nil {
		return nil, err
	}
	files = append(files, &GeneratedFile{
		Path:    filepath.Join("typed", "jzero", "credential.go"),
		Content: *bytes.NewBuffer(resourceGoBytes),
	})

	return files, nil
}
