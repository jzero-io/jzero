package ivminit

import (
	"github.com/jzero-io/jzero/cmd/jzero/internal/command/gen/gen"
	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

func (ivm *IvmInit) gen() error {
	defer gen.RemoveExtraFiles(config.C.Wd(), config.C.Ivm.Init.Style)

	err := ivm.jzeroRpc.Gen()
	if err != nil {
		return err
	}

	return nil
}
