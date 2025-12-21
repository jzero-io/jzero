package templatex

import (
	"strings"

	"github.com/jzero-io/jzero/core/templatex"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

// ParseTemplate template
func ParseTemplate(name string, data map[string]any, tplT []byte) ([]byte, error) {
	for _, v := range config.C.RegisterTplVal {
		split := strings.Split(v, "=")
		if len(split) == 2 {
			data[split[0]] = split[1]
		}
	}
	return templatex.ParseTemplateWithName(name, data, tplT, templatex.WithFuncMaps(registerFuncMap))
}
