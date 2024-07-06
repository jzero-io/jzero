package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

func getProtoDir(protoDirPath string) ([]os.DirEntry, error) {
	protoDir, err := os.ReadDir(protoDirPath)
	if err != nil {
		return nil, nil
	}
	return protoDir, nil
}

func GetProtoFilepath(protoDirPath string) ([]string, error) {
	var protoFilenames []string

	protoDir, err := getProtoDir(protoDirPath)
	if err != nil {
		return nil, err
	}

	for _, protoFile := range protoDir {
		if protoFile.IsDir() {
			filenames, err := GetProtoFilepath(filepath.Join(protoDirPath, protoFile.Name()))
			if err != nil {
				return nil, err
			}
			protoFilenames = append(protoFilenames, filenames...)
		} else {
			if strings.HasSuffix(protoFile.Name(), ".proto") {
				if b, err := protoHasService(filepath.Join(protoDirPath, protoFile.Name())); err == nil && b {
					protoFilenames = append(protoFilenames, filepath.Join(protoDirPath, protoFile.Name()))
				} else if err != nil {
					return nil, err
				}
			}
		}
	}
	return protoFilenames, nil
}

func protoHasService(fp string) (bool, error) {
	protoBaseDir := filepath.Join("desc", "proto")

	var protoParser protoparse.Parser
	protoParser.ImportPaths = []string{protoBaseDir}
	rel, err := filepath.Rel(protoBaseDir, fp)
	if err != nil {
		return false, err
	}
	fds, err := protoParser.ParseFiles(rel)
	if err != nil {
		return false, err
	}
	if len(fds) == 1 {
		if len(fds[0].GetServices()) >= 1 {
			return true, nil
		}
	}
	return false, nil
}

func GetMainApiFilePath(apiDirName string) (string, bool, error) {
	apiDir, err := os.ReadDir(apiDirName)
	if err != nil {
		return "", false, err
	}

	var mainApiFilePath string
	var isDelete bool

	for _, file := range apiDir {
		if file.Name() == "main.api" {
			mainApiFilePath = filepath.Join(apiDirName, file.Name())
			isDelete = false
			break
		}
	}

	if mainApiFilePath == "" {
		apiFilePath, err := getRouteApiFilePath(apiDirName)
		if err != nil {
			return "", false, err
		}
		sb := strings.Builder{}
		sb.WriteString("syntax = \"v1\"")
		sb.WriteString("\n")

		for _, api := range apiFilePath {
			sb.WriteString(fmt.Sprintf("import \"%s\"\n", api))
		}

		f, err := os.CreateTemp(apiDirName, "*.api")
		if err != nil {
			return "", false, err
		}

		_, err = f.WriteString(sb.String())
		if err != nil {
			return f.Name(), true, err
		}
		mainApiFilePath = f.Name()
		isDelete = true
		f.Close()
	}
	return mainApiFilePath, isDelete, nil
}

func GetApiServiceName(apiDirName string) string {
	fs, err := getRouteApiFilePath(apiDirName)
	if err != nil {
		return ""
	}
	for _, file := range fs {
		apiSpec, err := parser.Parse(filepath.Join(apiDirName, file), "")
		if err != nil {
			return ""
		}
		if apiSpec.Service.Name != "" {
			return apiSpec.Service.Name
		}
	}
	return ""
}
