package gensdk

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/generator"

	"github.com/jzero-io/jzero/embeded"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Scope        string
	ApiDir       string
	ProtoDir     string
	WrapResponse bool
	Dir          string
	WorkingDir   string
	Language     string
	Module       string

	Version string
)

func GenSdk(_ *cobra.Command, _ []string) error {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	if embeded.Home == "" {
		embeded.Home = filepath.Join(homeDir, ".jzero", Version)
	}

	// change dir
	if WorkingDir != "" {
		err := os.Chdir(WorkingDir)
		cobra.CheckErr(err)
	}

	if Dir != "" {
		if !pathx.FileExists(Dir) {
			if err := os.MkdirAll(Dir, 0o755); err != nil {
				cobra.CheckErr(err)
			}
		}
	}

	c := config.Config{
		Language:     Language,
		APP:          Scope,
		Module:       Module,
		Dir:          Dir,
		ApiDir:       ApiDir,
		ProtoDir:     ProtoDir,
		WrapResponse: WrapResponse,
	}

	gen, err := generator.New(c)
	if err != nil {
		return err
	}

	files, err := gen.Gen()
	if err != nil {
		return err
	}

	for _, v := range files {
		if !pathx.FileExists(filepath.Dir(filepath.Join(Dir, v.Path))) {
			if err = os.MkdirAll(filepath.Dir(filepath.Join(Dir, v.Path)), 0o755); err != nil {
				return err
			}
		}
		if pathx.FileExists(filepath.Join(Dir, v.Path)) && v.Skip {
			continue
		}
		if err = os.WriteFile(filepath.Join(Dir, v.Path), v.Content.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
