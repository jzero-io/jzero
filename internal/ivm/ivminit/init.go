package ivminit

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"path/filepath"
	"strings"
)

var (
	Version      string // version must v1 v2 v3...
	Style        string
	RemoveSuffix bool
)

type IvmInit struct {
	oldVersion   string
	newVersion   string
	protoBaseDir string
	oldProtoDir  string
	newProtoDir  string

	jzeroRpc gen.JzeroRpc
}

func Init(command *cobra.Command, args []string) error {
	var ivmInit IvmInit

	err := ivmInit.setOldVersion(Version)
	if err != nil {
		return err
	}
	ivmInit.newVersion = Version

	protoDir := filepath.Join("desc", "proto", ivmInit.oldVersion)
	protoBaseDir := filepath.Join("desc", "proto")
	ivmInit.protoBaseDir = protoBaseDir
	ivmInit.oldProtoDir = protoDir
	ivmInit.newProtoDir = filepath.Join("desc", "proto", Version)

	var protoFiles []string

	if pathx.FileExists(protoDir) {
		protoFiles, err = gen.GetProtoFilepath(protoDir)
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

	err = ivmInit.gen()
	if err != nil {
		return err
	}

	// invoke old version logic
	newVersionProtoFilepath, err := gen.GetProtoFilepath(ivmInit.newProtoDir)
	if err != nil {
		return err
	}

	for _, fp := range newVersionProtoFilepath {
		err = ivmInit.setUpdateProtoLogic(fp)
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
