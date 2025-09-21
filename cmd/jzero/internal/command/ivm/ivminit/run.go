package ivminit

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/genrpc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
	jzerodesc "github.com/jzero-io/jzero/cmd/jzero/internal/desc"
	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
)

type IvmInit struct {
	oldVersion   string
	newVersion   string
	protoBaseDir string
	oldProtoDir  string
	newProtoDir  string

	jzeroRpc genrpc.JzeroRpc
}

func Run(ic config.IvmConfig) error {
	var ivmInit IvmInit

	err := ivmInit.setOldVersion(ic.Version)
	if err != nil {
		return err
	}
	ivmInit.newVersion = ic.Version

	protoDir := filepath.Join("desc", "proto", ivmInit.oldVersion)
	protoBaseDir := filepath.Join("desc", "proto", "third_party")
	ivmInit.protoBaseDir = protoBaseDir
	ivmInit.oldProtoDir = protoDir
	ivmInit.newProtoDir = filepath.Join("desc", "proto", ic.Version)

	var protoFiles []string

	if pathx.FileExists(protoDir) {
		protoFiles, err = jzerodesc.GetProtoFilepath(protoDir)
		if err != nil {
			return err
		}
	}

	var fds []*desc.FileDescriptor

	// parse proto
	var protoParser protoparse.Parser

	protoParser.InferImportPaths = false

	if len(protoFiles) > 0 {
		protoParser.ImportPaths = []string{protoBaseDir}
		protoParser.IncludeSourceCodeInfo = true

		for _, protoFile := range protoFiles {
			rel, err := filepath.Rel(protoBaseDir, protoFile)
			if err != nil {
				return err
			}
			fds, err = protoParser.ParseFiles(rel)
			if err != nil {
				return err
			}
			for _, fd := range fds {
				err = ivmInit.updateProtoVersion(protoFile, fd)
				if err != nil {
					return err
				}
			}
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	moduleStruct, err := mod.GetGoMod(wd)
	if err != nil {
		return err
	}

	jzeroRpc := genrpc.JzeroRpc{
		Module: moduleStruct.Path,
	}
	ivmInit.jzeroRpc = jzeroRpc

	if config.C.Gen.Style == "" {
		config.C.Gen.Style = "gozero"
	}
	if config.C.Gen.Home == "" {
		config.C.Gen.Home = filepath.Join(config.C.Wd(), ".template")
	}

	err = ivmInit.gen()
	if err != nil {
		return err
	}

	// invoke old version logic
	newVersionProtoFilepath, err := jzerodesc.GetProtoFilepath(ivmInit.newProtoDir)
	if err != nil {
		return err
	}

	for i, fp := range newVersionProtoFilepath {
		oldFp := protoFiles[i]
		err = ivmInit.updateProtoLogic(fp, oldFp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ivm *IvmInit) setOldVersion(version string) error {
	v := cast.ToInt(strings.TrimPrefix(version, "v"))
	if v == 0 {
		return errors.New("please set version")
	}

	if v == 1 {
		return errors.New("version is v1, no need to init")
	}

	if v > 1 {
		ivm.oldVersion = "v" + cast.ToString(v-1)
		return nil
	}

	return nil
}
