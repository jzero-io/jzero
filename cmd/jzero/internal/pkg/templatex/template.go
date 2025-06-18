package templatex

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// ParseTemplate template
func ParseTemplate(name string, data any, tplT []byte) ([]byte, error) {
	var err error
	t := template.New(name).Funcs(sprig.TxtFuncMap())
	t.Funcs(registerFuncMap)

	t, err = t.Parse(string(tplT))
	if err != nil {
		return nil, err
	}

	t.Funcs(registerFuncMap)

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}
