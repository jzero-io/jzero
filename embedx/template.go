package embedx

import (
	"embed"
	"log"
	"os"
	"path/filepath"
)

var Template embed.FS

func ReadTemplateFile(filename string) []byte {
	data, _ := Template.ReadFile(filepath.Join(".template", filename))
	return data
}

func WriteTemplateDir(sourceDir, targetDir string) error {
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = writeTemplateDirRecursive(sourceDir, targetDir)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func writeTemplateDirRecursive(sourceDir, targetDir string) error {
	entries, err := Template.ReadDir(filepath.Join(".template", sourceDir))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return err
			}

			err = writeTemplateDirRecursive(sourcePath, targetPath)
			if err != nil {
				return err
			}
		} else {
			data, err := Template.ReadFile(filepath.Join(".template", sourcePath))
			if err != nil {
				return err
			}

			err = os.WriteFile(targetPath, data, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
