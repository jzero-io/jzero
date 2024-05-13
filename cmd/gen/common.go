package gen

import (
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

func GetProtoDir(wd string) ([]os.DirEntry, error) {
	protoDir, err := os.ReadDir(filepath.Join(wd, "app", "desc", "proto"))
	if err != nil {
		return nil, nil
	}
	return protoDir, nil
}

func GetProtoFilenames(wd string) ([]string, error) {
	protoDir, err := GetProtoDir(wd)
	if err != nil {
		return nil, nil
	}

	var protoFilenames []string
	for _, protoFile := range protoDir {
		if protoFile.IsDir() {
			continue
		}
		protoFilenames = append(protoFilenames, protoFile.Name())
	}
	return protoFilenames, nil
}

func GetMainApiFilePath(apiDirName string) string {
	apiDir, err := os.ReadDir(apiDirName)
	if err != nil {
		return ""
	}

	for _, file := range apiDir {
		if file.IsDir() {
			return GetMainApiFilePath(filepath.Join(apiDirName, file.Name()))
		} else {
			apiParser := parser.New(filepath.Join(apiDirName, file.Name()), "")
			apiAst := apiParser.Parse()
			for _, v := range apiAst.Stmts {
				switch v.(type) {
				case *ast.ImportGroupStmt, *ast.ImportLiteralStmt:
					return filepath.Join(apiDirName, file.Name())
				}
			}
		}
	}
	return ""
}
