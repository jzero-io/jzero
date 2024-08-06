package gensdk

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/config"
	gensdkconfig "github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/generator"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func GenSdk(gsc config.GenSdkConfig, genModule bool) error {
	if !pathx.FileExists(gsc.Output) {
		if err := os.MkdirAll(gsc.Output, 0o755); err != nil {
			return err
		}
	}

	c := gensdkconfig.Config{
		Language:     gsc.Language,
		Scope:        gsc.Scope,
		GenModule:    genModule,
		GoModule:     gsc.GoModule,
		GoPackage:    gsc.GoPackage,
		Output:       gsc.Output,
		ApiDir:       gsc.ApiDir,
		ProtoDir:     gsc.ProtoDir,
		WrapResponse: gsc.WrapResponse,
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
		if !pathx.FileExists(filepath.Dir(filepath.Join(gsc.Output, v.Path))) {
			if err = os.MkdirAll(filepath.Dir(filepath.Join(gsc.Output, v.Path)), 0o755); err != nil {
				return err
			}
		}
		if pathx.FileExists(filepath.Join(gsc.Output, v.Path)) && v.Skip {
			continue
		}
		if err = os.WriteFile(filepath.Join(gsc.Output, v.Path), v.Content.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
