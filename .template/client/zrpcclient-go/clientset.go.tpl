package {{ .APP }}

import (
	{{range $v := .Scopes}}"{{$.Module}}/typed/{{$v}}"{{end}}
)

type Interface interface {
	{{range $v := .Scopes}}{{$v | FirstUpper}}() {{$v}}.Interface{{end}}
}

type Clientset struct {
	{{range $v := .Scopes}}{{$v}} *{{$v}}.Client{{end}}
}

{{range $v := .Scopes}}func (x *Clientset) {{$v | FirstUpper}}() {{$v}}.Interface {
	return x.{{$v}}
}
{{ end }}

func NewClientset(ops ...Opt) *Clientset {
	cs := &Clientset{}

	for _, op := range ops {
		op(cs)
	}

	return cs
}
