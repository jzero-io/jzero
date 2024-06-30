package ivminit

import (
	"github.com/jzero-io/jzero/internal/gen"
	"github.com/jzero-io/jzero/pkg/mod"
	"os"
)

func (ivm *IvmInit) gen() error {
	defer gen.RemoveExtraFiles(ivm.jzeroRpc.Wd)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	moduleStruct, err := mod.GetGoMod(wd)
	if err != nil {
		return err
	}

	jzeroRpc := gen.JzeroRpc{
		Wd:           wd,
		Module:       moduleStruct.Path,
		Style:        Style,
		RemoveSuffix: RemoveSuffix,
	}
	ivm.jzeroRpc = jzeroRpc

	err = jzeroRpc.Gen()
	if err != nil {
		return err
	}

	return nil
}
