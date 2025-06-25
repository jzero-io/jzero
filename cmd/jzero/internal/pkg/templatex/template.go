package templatex

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jzero-io/jzero/cmd/jzero/internal/config"
)

// ParseTemplate template
func ParseTemplate(name string, data map[string]any, tplT []byte) ([]byte, error) {
	var err error
	t := template.New(name).Funcs(sprig.TxtFuncMap())
	t.Funcs(registerFuncMap)

	t, err = t.Parse(string(tplT))
	if err != nil {
		return nil, err
	}

	t.Funcs(registerFuncMap)

	buf := new(bytes.Buffer)

	logx.Debugf("get register tpl val: %v", config.C.RegisterTplVal)

	for _, v := range config.C.RegisterTplVal {
		split := strings.Split(v, "=")
		if len(split) == 2 {
			data[split[0]] = split[1]
		}
	}

	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}
