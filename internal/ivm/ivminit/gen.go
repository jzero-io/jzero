package ivminit

import (
	"github.com/jzero-io/jzero/internal/gen"
)

func (ivm *IvmInit) gen() error {
	defer gen.RemoveExtraFiles(ivm.jzeroRpc.Wd, ivm.jzeroRpc.Style)

	err := ivm.jzeroRpc.Gen()
	if err != nil {
		return err
	}

	return nil
}
