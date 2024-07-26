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
	Output       string
	Language     string
	GoModule     string
	GoPackage    string

	Version   string
	GenModule bool
)

func GenSdk(_ *cobra.Command, _ []string) error {
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	if embeded.Home == "" {
		embeded.Home = filepath.Join(homeDir, ".jzero", Version)
	}

	if Output != "" {
		if !pathx.FileExists(Output) {
			if err := os.MkdirAll(Output, 0o755); err != nil {
				cobra.CheckErr(err)
			}
		}
	}

	c := config.Config{
		Language: Language,
		Scope:    Scope,

		// 是否生成 go.mod 文件
		GenModule:    GenModule,
		GoModule:     GoModule,
		GoPackage:    GoPackage,
		Output:       Output,
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
		if !pathx.FileExists(filepath.Dir(filepath.Join(Output, v.Path))) {
			if err = os.MkdirAll(filepath.Dir(filepath.Join(Output, v.Path)), 0o755); err != nil {
				return err
			}
		}
		if pathx.FileExists(filepath.Join(Output, v.Path)) && v.Skip {
			continue
		}
		if err = os.WriteFile(filepath.Join(Output, v.Path), v.Content.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
