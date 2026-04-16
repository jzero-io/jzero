package gen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rinchsan/gosimports"
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
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console/progress"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
)

func Run() error {
	// 兼容之前的 gen style
	if config.C.Gen.Style != "" && config.C.Gen.Style != "gozero" {
		config.C.Style = config.C.Gen.Style
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

	jzeroModel := genmodel.JzeroModel{
		Module: module,
	}

	modelTitleFn := func() string {
		modelTitle := "model"
		if config.C.Gen.ModelDatasource {
			dsString := config.C.Gen.ModelDatasourceUrl[0]
			if len(config.C.Gen.ModelDatasourceUrl) > 1 {
				dsString = fmt.Sprintf("%s...", config.C.Gen.ModelDatasourceUrl[0])
			}
			modelTitle += " " + console.Cyan(fmt.Sprintf("by Datasource(%s)", dsString))
		}
		return console.Green("Gen") + " " + console.Yellow(modelTitle)
	}

	modelHeaderShown := config.C.Gen.GitChange && !config.C.Quiet && pathx.FileExists(config.C.SqlDir())

	// Show box header immediately for git-change mode only if sql dir exists
	if modelHeaderShown {
		modelTitle := "model"
		if config.C.Gen.ModelDatasource {
			modelTitle += " " + console.Cyan(fmt.Sprintf("by Datasource(%s)", strings.Join(config.C.Gen.ModelDatasourceUrl, ",")))
		}
		title := console.Green("Gen") + " " + console.Yellow(modelTitle) + " " + console.Cyan("(git-change mode)")
		fmt.Printf("%s\n", console.BoxHeader("", title))
	}

	// Generate model
	progressChan := make(chan progress.Message, 10)
	done := make(chan struct{})
	var modelErr error
	go func() {
		modelFiles, genErr := jzeroModel.Gen(progressChan)
		if genErr != nil {
			modelErr = genErr
		}
		_ = modelFiles
		close(done)
	}()

	modelState := progress.ConsumeStage(progressChan, done, modelTitleFn(), config.C.Quiet, modelHeaderShown)
	progress.FinishStage(modelTitleFn(), config.C.Quiet, &modelState, modelErr)

	if modelErr != nil {
		return modelErr
	}

	var apiSpecMap map[string]*spec.ApiSpec
	var protoSpecMap map[string]*rpcparser.Proto

	jzeroApi := genapi.JzeroApi{
		Module: module,
	}

	apiTitleFn := func() string {
		title := console.Green("Gen") + " " + console.Yellow("api")
		if config.C.Gen.GitChange {
			title += " " + console.Cyan("(git-change mode)")
		}
		return title
	}

	apiHeaderShown := !config.C.Quiet && pathx.FileExists(config.C.ApiDir())

	// Generate api
	apiProgressChan := make(chan progress.Message, 10)
	apiDone := make(chan struct{})
	var apiErr error

	// Show box header immediately before starting goroutine
	if apiHeaderShown {
		fmt.Printf("%s\n", console.BoxHeader("", apiTitleFn()))
	}

	go func() {
		apiSpecMap, apiErr = jzeroApi.Gen(apiProgressChan)
		close(apiDone)
	}()

	apiState := progress.ConsumeStage(apiProgressChan, apiDone, apiTitleFn(), config.C.Quiet, apiHeaderShown)
	progress.FinishStage(apiTitleFn(), config.C.Quiet, &apiState, apiErr)

	if apiErr != nil {
		return apiErr
	}

	jzeroRpc := genrpc.JzeroRpc{
		Module: module,
	}

	rpcTitleFn := func() string {
		title := console.Green("Gen") + " " + console.Yellow("rpc")
		if config.C.Gen.GitChange {
			title += " " + console.Cyan("(git-change mode)")
		}
		return title
	}

	rpcHeaderShown := config.C.Gen.GitChange && !config.C.Quiet && pathx.FileExists(config.C.ProtoDir())

	// Show box header immediately for git-change mode
	if rpcHeaderShown {
		fmt.Printf("%s\n", console.BoxHeader("", rpcTitleFn()))
	}

	// Generate rpc
	rpcProgressChan := make(chan progress.Message, 10)
	rpcDone := make(chan struct{})
	var rpcErr error
	go func() {
		protoSpecMap, rpcErr = jzeroRpc.Gen(rpcProgressChan)
		close(rpcProgressChan)
		close(rpcDone)
	}()

	rpcState := progress.ConsumeStage(rpcProgressChan, rpcDone, rpcTitleFn(), config.C.Quiet, rpcHeaderShown)
	progress.FinishStage(rpcTitleFn(), config.C.Quiet, &rpcState, rpcErr)

	if rpcErr != nil {
		return rpcErr
	}

	jzeroMongo := genmongo.JzeroMongo{
		Module: module,
	}

	// Only show mongo box if there are mongo types to generate
	var mongoShown bool
	if len(config.C.Gen.MongoType) > 0 {
		if !config.C.Quiet {
			title := console.Green("Gen") + " " + console.Yellow("mongo")
			fmt.Printf("%s\n", console.BoxHeader("", title))
			mongoShown = true
		}
		err = jzeroMongo.Gen()
		if !config.C.Quiet && mongoShown {
			if err != nil {
				fmt.Printf("%s\n\n", console.BoxErrorFooter())
			} else {
				fmt.Printf("%s\n\n", console.BoxSuccessFooter())
			}
		}
		if err != nil {
			return err
		}
	}

	// 收集并保存元数据（复用已解析的数据）
	if err = collectAndSaveMetadata(apiSpecMap, protoSpecMap); err != nil {
		// Debug removed("collect and save metadata error: %s", err.Error())
	}

	return nil
}

// collectAndSaveMetadata 收集并保存项目元数据（复用已解析的数据）
func collectAndSaveMetadata(apiSpecMap map[string]*spec.ApiSpec, protoSpecMap map[string]*rpcparser.Proto) error {
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
					_ = os.Remove(filepath.Join(wd, desc.GetApiFrameMainGoFilename(wd, v, style)))
				}
				if desc.GetApiFrameEtcFilename(wd, v, style) != "etc.yaml" {
					_ = os.Remove(filepath.Join(wd, "etc", desc.GetApiFrameEtcFilename(wd, v, style)))
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
					_ = os.Remove(filepath.Join(wd, desc.GetProtoFrameMainGoFilename(fileBase, style)))
				}
				if desc.GetProtoFrameEtcFilename(fileBase, style) != "etc.yaml" {
					_ = os.Remove(filepath.Join(wd, "etc", desc.GetProtoFrameEtcFilename(fileBase, style)))
				}
			}
		}
	}

	// Check if etc directory is empty and remove it if so
	etcDir := filepath.Join(wd, "etc")
	if pathx.FileExists(etcDir) {
		entries, err := os.ReadDir(etcDir)
		if err == nil && len(entries) == 0 {
			if err := os.Remove(etcDir); err != nil && !errors.Is(err, os.ErrNotExist) {
			}

		}
	}
}
