package embedx

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Opts func(config *embedxConfig)

type embedxConfig struct {
	Dir           string
	FileMatchFunc func(path string) bool
}

func WriteToLocal(ef embed.FS, opts ...Opts) ([]string, error) {
	config := &embedxConfig{}

	for _, opt := range opts {
		opt(config)
	}

	var fileList []string

	err := fs.WalkDir(ef, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			data, err := ef.ReadFile(path)
			if err != nil {
				return err
			}
			if config.Dir != "" {
				if stat, err := os.Stat(config.Dir); err != nil {
					if !os.IsExist(err) {
						err = os.MkdirAll(config.Dir, 0o755)
						if err != nil {
							return err
						}
					}
				} else {
					if !stat.IsDir() {
						return errors.Errorf("%s: not a directory", config.Dir)
					}
				}
			}

			var tmpFile *os.File
			if config.FileMatchFunc != nil {
				if config.FileMatchFunc(path) {
					if tmpFile, err = createTemp(config.Dir, path, data); err != nil {
						return err
					}
				}
			} else {
				if tmpFile, err = createTemp(config.Dir, path, data); err != nil {
					return err
				}
			}
			if tmpFile != nil {
				fileList = append(fileList, tmpFile.Name())
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}

func createTemp(dir, path string, data []byte) (*os.File, error) {
	tmpFile, err := os.CreateTemp(dir, fmt.Sprintf("*%s", filepath.Ext(path)))
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()
	if _, err := tmpFile.Write(data); err != nil {
		return nil, err
	}

	return tmpFile, nil
}
