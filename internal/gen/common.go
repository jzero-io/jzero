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
		if strings.HasSuffix(protoFile.Name(), ".proto") {
			protoFilenames = append(protoFilenames, protoFile.Name())
		}
	}
	return protoFilenames, nil
}

func GetMainApiFilePath(apiDirName string) (string, error) {
	apiDir, err := os.ReadDir(apiDirName)
	if err != nil {
		return "", err
	}

	var mainApiFilePath string

	for _, file := range apiDir {
		if file.Name() == "main.api" {
			mainApiFilePath = filepath.Join(apiDirName, file.Name())
			break
		}
	}

	if mainApiFilePath == "" {
		apiFilePath, err := getRouteApiFilePath(apiDirName)
		if err != nil {
			return "", err
		}
		sb := strings.Builder{}
		sb.WriteString("syntax = \"v1\"")
		sb.WriteString("\n")

		for _, api := range apiFilePath {
			sb.WriteString(fmt.Sprintf("import \"%s\"\n", api))
		}

		f, err := os.CreateTemp(apiDirName, "*.api")
		if err != nil {
			return "", err
		}

		_, err = f.WriteString(sb.String())
		if err != nil {
			return "", err
		}
		mainApiFilePath = f.Name()
		f.Close()
	}
	return mainApiFilePath, nil
}

func GetApiServiceName(apiDirName string) string {
	fs, err := getRouteApiFilePath(apiDirName)
	if err != nil {
		cobra.CheckErr(err)
	}
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