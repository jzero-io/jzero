package embeded

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Template embed.FS
	// Home template home
	Home string
)

func ReadTemplateFile(filename string) []byte {
	path := filepath.ToSlash(filepath.Join(".template", filename))
	data, err := Template.ReadFile(path)
	if err != nil {
		// 默认模板没有此文件, 但是用户自己的模板有这个文件
		if Home != "" {
			file, err := os.ReadFile(filepath.Join(Home, filename))
			if err == nil {
				return file
			}
		}
		return nil
	}
	if Home != "" {
		file, err := os.ReadFile(filepath.Join(Home, filename))
		if err != nil {
			return data
		}
		return file
	}
	return data
}

func ReadTemplateDir(dirname string) []os.DirEntry {
	path := filepath.ToSlash(filepath.Join(".template", dirname))
	data, err := Template.ReadDir(path)
	if err != nil {
		return nil
	}
	if Home != "" {
		file, err := os.ReadDir(filepath.Join(Home, dirname))
		if err != nil {
			return data
		}
		data = file
	}
	return data
}

func WriteTemplateDir(sourceDir, targetDir string) error {
	err := os.MkdirAll(targetDir, 0o755)
	if err != nil {
		return err
	}

	err = writeTemplateDirRecursive(sourceDir, targetDir)
	if err != nil {
		return err
	}

	return nil
}

func writeTemplateDirRecursive(sourceDir, targetDir string) error {
	entries, err := Template.ReadDir(filepath.ToSlash(filepath.Join(".template", sourceDir)))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			err := os.MkdirAll(targetPath, 0o755)
			if err != nil {
				return err
			}

			err = writeTemplateDirRecursive(sourcePath, targetPath)
			if err != nil {
				return err
			}
		} else {
			if pathx.FileExists(targetPath) {
				continue
			}

			data, err := Template.ReadFile(filepath.ToSlash(filepath.Join(".template", sourcePath)))
			if err != nil {
				return err
			}

			file, err := os.ReadFile(filepath.Join(Home, sourcePath))
			if err == nil {
				data = file
			}

			err = os.WriteFile(targetPath, data, 0o644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
