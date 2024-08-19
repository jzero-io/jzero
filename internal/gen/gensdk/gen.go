package gensdk

import (
	"os"
	"path/filepath"

	"github.com/jzero-io/jzero/config"
	gensdkconfig "github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/generator"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func GenSdk(gc config.GenConfig, genModule bool) error {
	if !pathx.FileExists(gc.Sdk.Output) {
		if err := os.MkdirAll(gc.Sdk.Output, 0o755); err != nil {
			return err
		}
	}

	c := gensdkconfig.Config{
		Language:     gc.Sdk.Language,
		Scope:        gc.Sdk.Scope,
		GenModule:    genModule,
		GoVersion:    gc.Sdk.GoVersion,
		GoModule:     gc.Sdk.GoModule,
		GoPackage:    gc.Sdk.GoPackage,
		Output:       gc.Sdk.Output,
		ApiDir:       gc.Sdk.ApiDir,
		ProtoDir:     gc.Sdk.ProtoDir,
		WrapResponse: gc.Sdk.WrapResponse,
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
		if !pathx.FileExists(filepath.Dir(filepath.Join(gc.Sdk.Output, v.Path))) {
			if err = os.MkdirAll(filepath.Dir(filepath.Join(gc.Sdk.Output, v.Path)), 0o755); err != nil {
				return err
			}
		}
		if pathx.FileExists(filepath.Join(gc.Sdk.Output, v.Path)) && v.Skip {
			continue
		}
		if err = os.WriteFile(filepath.Join(gc.Sdk.Output, v.Path), v.Content.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
