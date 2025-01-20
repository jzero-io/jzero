package gensdk

import (
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	gensdkconfig "github.com/jzero-io/jzero/internal/gen/gensdk/config"
	"github.com/jzero-io/jzero/internal/gen/gensdk/generator"
)

func GenSdk(genModule bool) error {
	if !pathx.FileExists(config.C.Gen.Sdk.Output) {
		if err := os.MkdirAll(config.C.Gen.Sdk.Output, 0o755); err != nil {
			return err
		}
	}

	gc := gensdkconfig.Config{
		GenModule: genModule,
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
		if !pathx.FileExists(filepath.Dir(filepath.Join(config.C.Gen.Sdk.Output, v.Path))) {
			if err = os.MkdirAll(filepath.Dir(filepath.Join(config.C.Gen.Sdk.Output, v.Path)), 0o755); err != nil {
				return err
			}
		}
		if pathx.FileExists(filepath.Join(config.C.Gen.Sdk.Output, v.Path)) && v.Skip {
			continue
		}
		if err = os.WriteFile(filepath.Join(config.C.Gen.Sdk.Output, v.Path), v.Content.Bytes(), 0o644); err != nil {
			return err
		}
	}

	return nil
}
