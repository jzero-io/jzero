package plugins

import (
	"github.com/zeromicro/go-zero/zrpc"
{{range $serviceName := .Services}}
	{{$.PluginName | ToCamel | FirstLower}}{{$serviceName | ToCamel | FirstUpper}} "{{$.Module}}/plugins/{{$.PluginName}}/typed/{{$serviceName | lower}}"{{end}}
)

type {{.PluginName | ToCamel | FirstUpper}} interface {
{{range $serviceName := .Services}}	{{$serviceName | ToCamel | FirstUpper}}() {{$.PluginName | ToCamel | FirstLower}}{{$serviceName | ToCamel | FirstUpper}}.{{$serviceName | ToCamel | FirstUpper}}
{{end}}
}

type {{.PluginName | ToCamel | FirstLower}}Client struct {
	conn zrpc.Client
{{range $serviceName := .Services}}	{{$serviceName | ToCamel | FirstLower}} {{$.PluginName | ToCamel | FirstLower}}{{$serviceName | ToCamel | FirstUpper}}.{{$serviceName | ToCamel | FirstUpper}}
{{end}}
}

{{range $serviceName := .Services}}func (x *{{$.PluginName | ToCamel | FirstLower}}Client) {{$serviceName | ToCamel | FirstUpper}}() {{$.PluginName | ToCamel | FirstLower}}{{$serviceName | ToCamel | FirstUpper}}.{{$serviceName | ToCamel | FirstUpper}} {
	return x.{{$serviceName | ToCamel | FirstLower}}
}

{{end}}

func New{{.PluginName | ToCamel | FirstUpper}}(conn zrpc.Client) {{.PluginName | ToCamel | FirstUpper}} {
	return &{{.PluginName | ToCamel | FirstLower}}Client{
		conn: conn,
{{range $serviceName := .Services}}		{{$serviceName | ToCamel | FirstLower}}: {{$.PluginName | ToCamel | FirstLower}}{{$serviceName | ToCamel | FirstUpper}}.New{{$serviceName | ToCamel | FirstUpper}}(conn),
{{end}}
	}
}