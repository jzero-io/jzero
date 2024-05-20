package templatex

import (
	"bytes"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/jzero-io/jzero/app/pkg/stringx"
)

// ParseTemplate template
func ParseTemplate(data interface{}, tplT []byte) ([]byte, error) {
	t := template.Must(template.New("production").Funcs(sprig.TxtFuncMap()).Funcs(RegisterTxtFuncMap()).Parse(string(tplT)))

	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func RegisterTxtFuncMap() template.FuncMap {
	return RegisterFuncMap()
}

func RegisterFuncMap() template.FuncMap {
	gfm := make(map[string]interface{}, len(registerFuncMap))
	for k, v := range registerFuncMap {
		gfm[k] = v
	}
	return gfm
}

var registerFuncMap = map[string]interface{}{
	"FirstUpper": stringx.FirstUpper,
	"ToCamel":    stringx.ToCamel,
}
