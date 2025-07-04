// Code generated by jzero. DO NOT EDIT.

package {{.Package}}

import (
    "github.com/zeromicro/go-zero/zrpc"

	{{range $v := .Services}}{{$v | ToCamel | lower}} "{{$.Module}}/typed/{{$v | lower}}"
	{{end}}
)

type Clientset interface {
	{{range $v := .Services}}{{$v | ToCamel | FirstUpper}}() {{$v | ToCamel | lower}}.{{$v | ToCamel | FirstUpper}}
	{{end}}}

type clientset struct {
	{{range $v := .Services}}{{$v | ToCamel | FirstLower}} {{$v | ToCamel | lower}}.{{$v | ToCamel | FirstUpper}}
	{{end}}}

{{range $v := .Services}}func (cs *clientset) {{$v | FirstUpper | ToCamel}}() {{$v | ToCamel |lower}}.{{$v | ToCamel | FirstUpper}} {
	return cs.{{$v | ToCamel | FirstLower}}
}

{{end}}

func NewClientset(cli zrpc.Client) (Clientset, error) {
    cs := clientset{
		{{range $v := .Services}}{{$v | ToCamel | FirstLower}}: {{$v | ToCamel | lower}}.New{{$v | ToCamel | FirstUpper}}(cli),
		{{end}}}

	return &cs, nil
}