package pb

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	//go:embed *
	Embed embed.FS
)

func WriteToLocal(ef embed.FS) ([]string, error) {
	var fileList []string

	err := fs.WalkDir(ef, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".pb" {
			data, err := ef.ReadFile(path)
			if err != nil {
				return err
			}
			pbPath := filepath.Join("desc", "pb", path)
			fileList = append(fileList, pbPath)
			if err := os.MkdirAll(filepath.Dir(pbPath), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(pbPath, data, 0o644); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}
