// Code generated by jzero. DO NOT EDIT.
// type: {{.Scope}}_client

package {{.Scope}}

import (
	"{{.Module}}/rest"
)

type {{.Scope | FirstUpper}}Interface interface {
	RESTClient() rest.Interface
	
	{{range $v := .Resources}}{{$v | FirstUpper}}Getter
	{{end}}
}

type {{.Scope | FirstUpper}}Client struct {
	restClient rest.Interface
}

func (x *{{.Scope | FirstUpper}}Client) RESTClient() rest.Interface {
	if x == nil {
		return nil
	}
	return x.restClient
}

{{range $v := .Resources}}func (x *{{$.Scope | FirstUpper}}Client) {{$v | FirstUpper}}() {{$v | FirstUpper}}Interface {
	return new{{$v | FirstUpper}}Client(x)
}

{{end}}
// NewForConfig creates a new {{.Scope | FirstUpper}}Client for the given config.
func NewForConfig(x *rest.RESTClient) (*{{.Scope | FirstUpper}}Client, error) {
	config := *x
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &{{.Scope | FirstUpper}}Client{client}, nil
}