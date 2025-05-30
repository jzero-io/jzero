package gensdk

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	gensdkconfig "github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gensdk/generator"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
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

		if filepath.Ext(v.Path) == ".go" {
			formated, err := gosimports.Process("", v.Content.Bytes(), nil)
			if err != nil {
				return errors.Errorf("format go file %s %s meet error: %v", v.Path, v.Content.Bytes(), err)
			}
			if err = os.WriteFile(filepath.Join(config.C.Gen.Sdk.Output, v.Path), formated, 0o644); err != nil {
				return err
			}
		} else {
			if err = os.WriteFile(filepath.Join(config.C.Gen.Sdk.Output, v.Path), v.Content.Bytes(), 0o644); err != nil {
				return err
			}
		}
	}

	return nil
}
