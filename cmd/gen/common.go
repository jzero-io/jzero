package gen

import (
	"os"
	"path/filepath"
)

func GetProtoDir(wd string) ([]os.DirEntry, error) {
	protoDir, err := os.ReadDir(filepath.Join(wd, "daemon", "desc", "proto"))
	if err != nil {
		return nil, err
	}
	return protoDir, nil
}

func GetProtoFilenames(wd string) ([]string, error) {
	protoDir, err := GetProtoDir(wd)
	if err != nil {
		return nil, err
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
