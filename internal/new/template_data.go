package new

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/jzero-io/jzero/pkg/mod"
)

const (
	go1_21_0 = "1.21.0"
)

func newTemplateData() (map[string]interface{}, error) {
	goVersion, err := mod.GetGoVersion()
	if err != nil {
		return nil, err
	}

	newVersion, err := version.NewVersion(goVersion)
	if err != nil {
		return nil, err
	}

	go1210version, err := version.NewVersion(go1_21_0)
	if err != nil {
		return nil, err
	}

	if newVersion.LessThan(go1210version) {
		split := strings.Split(goVersion, ".")
		goVersion = fmt.Sprintf("%s.%s", split[0], split[1])
	}

	templateData := map[string]interface{}{
		"Module":    Module,
		"APP":       AppName,
		"GoVersion": goVersion,
		"GoArch":    runtime.GOARCH,
	}

	return templateData, nil
}
