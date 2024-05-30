package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

func GetProtoDir(protoDirPath string) ([]os.DirEntry, error) {
	protoDir, err := os.ReadDir(protoDirPath)
	if err != nil {
		return nil, nil
	}
	return protoDir, nil
}

func GetProtoFilenames(protoDirPath string) ([]string, error) {
	protoDir, err := GetProtoDir(protoDirPath)
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

	var mainApiFilePath string

	for _, file := range apiDir {
		if file.Name() == "main.api" {
			mainApiFilePath = filepath.Join(apiDirName, file.Name())
			break
		}
	}

	if mainApiFilePath == "" {
		apiFilePath := getRouteApiFilePath(apiDirName)
		sb := strings.Builder{}
		sb.WriteString("syntax = \"v1\"")
		sb.WriteString("\n")

		for _, api := range apiFilePath {
			sb.WriteString(fmt.Sprintf("import \"%s\"\n", api))
		}

		f, err := os.CreateTemp(apiDirName, "*.api")
		if err != nil {
			return ""
		}

		_, err = f.WriteString(sb.String())
		if err != nil {
			return ""
		}
		mainApiFilePath = f.Name()
		f.Close()
	}
	return mainApiFilePath
}

func GetApiServiceName(apiDirName string) string {
	fs := getRouteApiFilePath(apiDirName)
	for _, file := range fs {
		apiSpec, err := parser.Parse(filepath.Join(apiDirName, file), "")
		if err != nil {
			cobra.CheckErr(err)
		}
		if apiSpec.Service.Name != "" {
			return apiSpec.Service.Name
		}
	}
	return ""
}
