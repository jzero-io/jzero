package gensdk

import (
	"github.com/jaronnie/jzero/cmd/gensdk/generator"
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
	if Dir != "" {
		if !pathx.FileExists(Dir) {
			if err := os.MkdirAll(Dir, 0o755); err != nil {
				return err
			}
		}
	}

	g, err := generator.New(Language)
	if err != nil {
		return err
	}

	files, err := g.Gen()
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
