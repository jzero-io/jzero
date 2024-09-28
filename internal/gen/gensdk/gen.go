package gensdk

import (
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	gensdkconfig "github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/generator"
)

func GenSdk(c config.Config, genModule bool) error {
	if !pathx.FileExists(c.Gen.Sdk.Output) {
		if err := os.MkdirAll(c.Gen.Sdk.Output, 0o755); err != nil {
			return err
		}
	}

	gc := gensdkconfig.Config{
		Language:     c.Gen.Sdk.Language,
		Scope:        c.Gen.Sdk.Scope,
		GenModule:    genModule,
		GoVersion:    c.Gen.Sdk.GoVersion,
		GoModule:     c.Gen.Sdk.GoModule,
		GoPackage:    c.Gen.Sdk.GoPackage,
		Output:       c.Gen.Sdk.Output,
		ApiDir:       c.Gen.Sdk.ApiDir,
		ProtoDir:     c.Gen.Sdk.ProtoDir,
		WrapResponse: c.Gen.Sdk.WrapResponse,
	}

	gen, err := generator.New(gc)
	if err != nil {
		return err
	}

	files, err := gen.Gen()
	if err != nil {
		return err
	}

	for _, v := range files {
		if !pathx.FileExists(filepath.Dir(filepath.Join(c.Gen.Sdk.Output, v.Path))) {
			if err = os.MkdirAll(filepath.Dir(filepath.Join(c.Gen.Sdk.Output, v.Path)), 0o755); err != nil {
				return err
			}
		}
		if pathx.FileExists(filepath.Join(c.Gen.Sdk.Output, v.Path)) && v.Skip {
			continue
		}
		if err = os.WriteFile(filepath.Join(c.Gen.Sdk.Output, v.Path), v.Content.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
