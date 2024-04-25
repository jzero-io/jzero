package generator

import (
	"bytes"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"path/filepath"

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

	_, err = jparser.Parse(fds, apiSpecs)
	if err != nil {
		return nil, err
	}

	var files []*GeneratedFile

	files = append(files, &GeneratedFile{
		Path:    "client.go",
		Content: *bytes.NewBuffer([]byte("package client")),
	})
	return files, nil
}
