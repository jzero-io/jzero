package plugin

import (
	"os"
	"path/filepath"

	ddlparser "github.com/zeromicro/ddl-parser/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"

	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
)

type Desc struct {
	Api   Api
	Proto Proto
	Model Model
}

type Api struct {
	SpecMap map[string]*spec.ApiSpec
}

type Proto struct {
	SpecMap map[string]*rpcparser.Proto
}

type Model struct {
	SpecMap map[string]*ddlparser.Table
}

type Plugin struct {
	Name string
	Desc Desc
}

type Metadata struct {
	Desc   Desc
	Plugin Plugin // serverless plugin
}

// New creates a new jzero cmd plugin with parsed API, Proto, and SQL specifications
func New() (*Metadata, error) {
	metadata := &Metadata{}

	// Parse API files if they exist
	apiSpecMap, err := parseApiFiles(filepath.Join("desc", "api"))
	if err == nil && len(apiSpecMap) > 0 {
		metadata.Desc.Api.SpecMap = apiSpecMap
	}

	// Parse Proto files if they exist
	protoSpecMap, err := parseProtoFiles(filepath.Join("desc", "proto"))
	if err == nil && len(protoSpecMap) > 0 {
		metadata.Desc.Proto.SpecMap = protoSpecMap
	}

	// Parse SQL files if they exist
	sqlSpecMap, err := parseSqlFiles(filepath.Join("desc", "sql"))
	if err == nil && len(sqlSpecMap) > 0 {
		metadata.Desc.Model.SpecMap = sqlSpecMap
	}

	return metadata, nil
}

// parseApiFiles parses all API files in the given directory
func parseApiFiles(dir string) (map[string]*spec.ApiSpec, error) {
	files, err := desc.FindApiFiles(dir)
	if err != nil || len(files) == 0 {
		return nil, err
	}

	specMap := make(map[string]*spec.ApiSpec)
	for _, file := range files {
		apiSpec, err := parser.Parse(file, "")
		if err != nil {
			return nil, err
		}
		specMap[file] = apiSpec
	}

	return specMap, nil
}

// parseProtoFiles parses all Proto files in the given directory
func parseProtoFiles(dir string) (map[string]*rpcparser.Proto, error) {
	files, err := desc.FindRpcServiceProtoFiles(dir)
	if err != nil || len(files) == 0 {
		return nil, err
	}

	specMap := make(map[string]*rpcparser.Proto)
	protoParser := rpcparser.NewDefaultProtoParser()

	for _, file := range files {
		var protoSpec rpcparser.Proto
		protoSpec, err = protoParser.Parse(file, true)
		if err != nil {
			return nil, err
		}
		specMap[file] = &protoSpec
	}

	return specMap, nil
}

// parseSqlFiles parses all SQL files in the given directory
func parseSqlFiles(dir string) (map[string]*ddlparser.Table, error) {
	files, err := desc.FindSqlFiles(dir)
	if err != nil || len(files) == 0 {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	specMap := make(map[string]*ddlparser.Table)
	for _, file := range files {
		p := ddlparser.NewParser()
		tables, err := p.From(filepath.Join(wd, file))
		if err != nil {
			return nil, err
		}

		for _, table := range tables {
			specMap[table.Name] = table
		}
	}

	return specMap, nil
}
