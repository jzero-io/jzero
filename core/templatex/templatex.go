package templatex

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/eddieowens/opts"
)

type TemplateOpts struct {
	Options  []string
	FuncMaps []template.FuncMap
}

func (opts TemplateOpts) DefaultOptions() TemplateOpts {
	return TemplateOpts{}
}

func WithFuncMaps(funcMaps []template.FuncMap) opts.Opt[TemplateOpts] {
	return func(o *TemplateOpts) {
		o.FuncMaps = funcMaps
	}
}

func WithOptions(options ...string) opts.Opt[TemplateOpts] {
	return func(o *TemplateOpts) {
		o.Options = options
	}
}

// ParseTemplate template
func ParseTemplate(data any, tplT []byte, op ...opts.Opt[TemplateOpts]) ([]byte, error) {
	o := opts.DefaultApply(op...)

	var err error
	t := template.New("production").Funcs(sprig.TxtFuncMap())
	for _, funcMap := range o.FuncMaps {
		t = t.Funcs(funcMap)
	}
	if len(o.Options) > 0 {
		t = t.Option(o.Options...)
	}

	t, err = t.Parse(string(tplT))
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
