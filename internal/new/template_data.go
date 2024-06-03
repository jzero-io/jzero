package new

import "github.com/jzero-io/jzero/pkg/mod"

func newTemplateData() (map[string]interface{}, error) {
	goVersion, err := mod.GetGoVersion()
	if err != nil {
		return nil, err
	}

	templateData := map[string]interface{}{
		"Module":    Module,
		"APP":       AppName,
		"GoVersion": goVersion,
	}

	return templateData, nil
}
