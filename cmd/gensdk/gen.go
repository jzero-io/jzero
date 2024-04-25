package gensdk

import (
	"github.com/jaronnie/genius"
	"github.com/jaronnie/jzero/cmd/gensdk/generator"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"os"
	"path/filepath"
)

var (
	Dir      string
	Language string
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

	target := generator.Target{
		Language: Language,
		APP:      cast.ToString(g.Get("APP")),
	}

	gen, err := generator.New(target)
	if err != nil {
		return err
	}

	files, err := gen.Gen()
	if err != nil {
		return err
	}

	for _, v := range files {
		if !pathx.FileExists(filepath.Dir(v.Path)) {
			if err = os.MkdirAll(filepath.Dir(v.Path), 0o755); err != nil {
				return err
			}
		}
		if err = os.WriteFile(filepath.Join(Dir, v.Path), v.Content.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
