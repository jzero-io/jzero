package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen/genapi"
	"github.com/jzero-io/jzero/internal/gen/genmodel"
	"github.com/jzero-io/jzero/internal/gen/genrpc"
	"github.com/jzero-io/jzero/pkg/desc"
	"github.com/jzero-io/jzero/pkg/mod"
)

func Run() error {
	fmt.Printf("%s working dir %s\n", color.WithColor("Enter", color.FgGreen), config.C.Wd())

	var module string
	moduleStruct, err := mod.GetGoMod(config.C.Wd())
	if err != nil {
		return errors.Wrapf(err, "get go module struct error")
	}
	module = moduleStruct.Path

	if !pathx.FileExists("go.mod") {
		module, err = mod.GetParentPackage(config.C.Wd())
		if err != nil {
			return errors.Wrapf(err, "get parent package error")
		}
	}

	defer func() {
		RemoveExtraFiles(config.C.Wd(), config.C.Gen.Style)
	}()

	jzeroModel := genmodel.JzeroModel{
		Module: module,
	}
	err = jzeroModel.Gen()
	if err != nil {
		return err
	}

	jzeroApi := genapi.JzeroApi{
		Module: module,
	}
	err = jzeroApi.Gen()
	if err != nil {
		return err
	}

	jzeroRpc := genrpc.JzeroRpc{
		Module: module,
	}
	err = jzeroRpc.Gen()
	if err != nil {
		return err
	}

	return nil
}

func RemoveExtraFiles(wd, style string) {
	if pathx.FileExists(filepath.Join("desc", "api")) {
		apiFilenames, err := desc.FindApiFiles(filepath.Join("desc", "api"))
		if err == nil {
			for _, v := range apiFilenames {
				if desc.GetApiFrameMainGoFilename(wd, v, style) != "main.go" {
					if err := os.Remove(filepath.Join(wd, desc.GetApiFrameMainGoFilename(wd, v, style))); err != nil && !errors.Is(err, os.ErrNotExist) {
						logx.Debugf("remove api frame main go file error: %s", err.Error())
					}
				}
				if desc.GetApiFrameEtcFilename(wd, v, style) != "etc.yaml" {
					if err := os.Remove(filepath.Join(wd, "etc", desc.GetApiFrameEtcFilename(wd, v, style))); err != nil && !errors.Is(err, os.ErrNotExist) {
						logx.Debugf("remove api etc file error: %s", err.Error())
					}
				}
			}
		}
	}

	if pathx.FileExists(filepath.Join("desc", "proto")) {
		protoFilenames, err := desc.GetProtoFilepath(filepath.Join("desc", "proto"))
		if err == nil {
			for _, v := range protoFilenames {
				v = filepath.Base(v)
				fileBase := v[0 : len(v)-len(path.Ext(v))]
				if desc.GetProtoFrameMainGoFilename(fileBase, style) != "main.go" {
					if err = os.Remove(filepath.Join(wd, desc.GetProtoFrameMainGoFilename(fileBase, style))); err != nil && !errors.Is(err, os.ErrNotExist) {
						logx.Debugf("remove proto frame main go file error: %s", err.Error())
					}
				}
				if desc.GetProtoFrameEtcFilename(fileBase, style) != "etc.yaml" {
					if err = os.Remove(filepath.Join(wd, "etc", desc.GetProtoFrameEtcFilename(fileBase, style))); err != nil && !errors.Is(err, os.ErrNotExist) {
						logx.Debugf("remove proto etc file error: %s", err.Error())
					}
				}
			}
		}
	}
}
