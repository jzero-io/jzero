package embeded

import (
	"embed"
	"os"
	"path/filepath"
)

var (
	Template embed.FS

	// Home template home
	Home string
)

func ReadTemplateFile(filename string) []byte {
	if Home != "" {
		file, err := os.ReadFile(filepath.Join(Home, filename))
		if err == nil {
			return file
		}
	}
	path := filepath.ToSlash(filepath.Join(".template", filename))
	data, err := Template.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func ReadTemplateDir(dirname string) []os.DirEntry {
	if Home != "" {
		file, err := os.ReadDir(filepath.Join(Home, dirname))
		if err == nil {
			return file
		}
	}
	path := filepath.ToSlash(filepath.Join(".template", dirname))
	data, err := Template.ReadDir(path)
	if err != nil {
		return nil
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
			data, err := Template.ReadFile(filepath.ToSlash(filepath.Join(".template", sourcePath)))
			if err != nil {
				return err
			}

			err = os.WriteFile(targetPath, data, 0o644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
