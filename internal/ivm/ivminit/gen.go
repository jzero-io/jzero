package ivminit

import (
	"github.com/jzero-io/jzero/config"
	"github.com/jzero-io/jzero/internal/gen"
)

func (ivm *IvmInit) gen() error {
	defer gen.RemoveExtraFiles(config.C.Wd(), config.C.Ivm.Init.Style)

	err := ivm.jzeroRpc.Gen()
	if err != nil {
		return err
	}

	return nil
}
