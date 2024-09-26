package templatex

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/jzero-io/jzero/pkg/stringx"
)

// ParseTemplate template
func ParseTemplate(data interface{}, tplT []byte) ([]byte, error) {
	t, err := template.New("production").Funcs(sprig.TxtFuncMap()).Funcs(RegisterTxtFuncMap()).Option().Parse(string(tplT))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
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
	"FirstLower": stringx.FirstLower,
	"ToCamel":    stringx.ToCamel,
}
