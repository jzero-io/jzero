package {{ .Package }}

import (
	{{range $v := .Scopes}}"{{$.Module}}/{{ $.ClientDir }}"{{end}}
)

type Interface interface {
	{{range $v := .Scopes}}{{$v | FirstUpper | ToCamel}}() {{$v | lower}}.Interface{{end}}
}

type Clientset struct {
	{{range $v := .Scopes}}{{$v | ToCamel}} *{{$v}}.Client{{end}}
}

{{range $v := .Scopes}}func (x *Clientset) {{$v | FirstUpper | ToCamel}}() {{$v | lower}}.Interface {
	return x.{{$v | ToCamel}}
}
{{ end }}

func NewClientWithOptions(ops ...Opt) *Clientset {
	cs := &Clientset{}

	for _, op := range ops {
		op(cs)
	}

	return cs
}
