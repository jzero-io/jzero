package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	rpcparser "github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genapi"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genmodel"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genmongo"
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genrpc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	"github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
)

type OperateGenT struct {
	GenModel bool
	GenApi   bool
	GenRpc   bool
	GenMongo bool
}

func Run() error {
	// 兼容之前的 gen style
	if config.C.Gen.Style != "" && config.C.Gen.Style != "gozero" {
		config.C.Style = config.C.Gen.Style
	}

	if !config.C.Quiet {
		fmt.Printf("%s working dir %s\n", console.Green("Enter"), config.C.Wd())
	}

	var module string
	moduleStruct, err := mod.GetGoMod(config.C.Wd())
	if err != nil {
		return errors.Wrapf(err, "get go module struct error")
	}
	module = moduleStruct.Path
	gosimports.LocalPrefix = module

	if !pathx.FileExists("go.mod") {
		module, err = mod.GetParentPackage(config.C.Wd())
		if err != nil {
			return errors.Wrapf(err, "get parent package error")
		}
	}

	defer func() {
		RemoveExtraFiles(config.C.Wd(), config.C.Style)
	}()

	// 当存在 desc 指定时，只生成指定文件 (暂不支持指定 mongo)
	var operateGenData OperateGenT
	if len(config.C.Gen.Desc) >= 1 {
		for _, v := range config.C.Gen.Desc {
			if strings.HasSuffix(v, ".api") {
				operateGenData.GenApi = true
			}
			if strings.HasSuffix(v, ".proto") {
				operateGenData.GenRpc = true
			}
			if strings.HasSuffix(v, ".sql") {
				operateGenData.GenModel = true
			}
		}
	} else {
		operateGenData.GenModel, operateGenData.GenApi, operateGenData.GenRpc, operateGenData.GenMongo = true, true, true, true
	}

	apiSpecMap, protoSpecMap, err := operateGen(module, operateGenData)
	if err != nil {
		return err
	}
	// 收集并保存元数据（复用已解析的数据）
	if err = collectAndSaveMetadata(apiSpecMap, protoSpecMap); err != nil {
		logx.Debugf("collect and save metadata error: %s", err.Error())
	}

	return nil
}

func operateGen(module string, t OperateGenT) (apiSpecMap map[string]*spec.ApiSpec, protoSpecMap map[string]rpcparser.Proto, err error) {
	if t.GenModel {
		jzeroModel := genmodel.JzeroModel{
			Module: module,
		}
		err = jzeroModel.Gen()
		if err != nil {
			return nil, nil, err
		}
	}

	if t.GenApi {
		jzeroApi := genapi.JzeroApi{
			Module: module,
		}
		apiSpecMap, err = jzeroApi.Gen()
		if err != nil {
			return nil, nil, err
		}
	}

	if t.GenRpc {
		jzeroRpc := genrpc.JzeroRpc{
			Module: module,
		}
		protoSpecMap, err = jzeroRpc.Gen()
		if err != nil {
			return nil, nil, err
		}
	}
	if t.GenMongo {
		jzeroMongo := genmongo.JzeroMongo{
			Module: module,
		}
		err = jzeroMongo.Gen()
		if err != nil {
			return nil, nil, err
		}
	}
	return
}

// collectAndSaveMetadata 收集并保存项目元数据（复用已解析的数据）
func collectAndSaveMetadata(apiSpecMap map[string]*spec.ApiSpec, protoSpecMap map[string]rpcparser.Proto) error {
	if len(apiSpecMap) == 0 && len(protoSpecMap) == 0 {
		return nil
	}

	var md desc.Metadata

	if len(apiSpecMap) > 0 {
		apiMetadata, err := desc.CollectFromAPI(apiSpecMap)
		if err != nil {
			return errors.Wrapf(err, "collect api metadata")
		}
		md.API = apiMetadata
	}

	if len(protoSpecMap) > 0 {
		protoMetadata, err := desc.CollectFromProto(protoSpecMap)
		if err != nil {
			return errors.Wrapf(err, "collect proto metadata")
		}
		md.Proto = protoMetadata
	}

	if err := desc.Save(&md); err != nil {
		return errors.Wrapf(err, "save metadata")
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
		protoFilenames, err := desc.FindRpcServiceProtoFiles(filepath.Join("desc", "proto"))
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
