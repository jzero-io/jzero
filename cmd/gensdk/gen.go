package gensdk

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/cmd/gensdk/config"

	"github.com/jaronnie/genius"
	"github.com/jzero-io/jzero/cmd/gensdk/generator"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	Dir      string
	Language string
	Module   string
)

func GenSdk(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	cobra.CheckErr(err)
	configBytes, err := os.ReadFile(filepath.Join(wd, "config.toml"))
	cobra.CheckErr(err)

	g, err := genius.NewFromToml(configBytes)
	cobra.CheckErr(err)

	if Dir != "" {
		if !pathx.FileExists(Dir) {
			if err := os.MkdirAll(Dir, 0o755); err != nil {
				return err
			}
		}
	}

	c := config.Config{
		Language: Language,
		APP:      cast.ToString(g.Get("APP")),
		Module:   Module,
		Dir:      Dir,
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
