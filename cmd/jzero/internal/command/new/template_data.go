package new

import (
	"runtime"

	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/mod"
)

func NewTemplateData() (map[string]any, error) {
	goVersion, err := mod.GetGoVersion()
	if err != nil {
		return nil, err
	}

	templateData := map[string]any{
		"GoVersion": goVersion,
		"GoArch":    runtime.GOARCH,
	}

	return templateData, nil
}
