package plugins

import (
	"github.com/jzero-io/jzero/core/restc"
)

type Plugins interface {
{{range $pluginName := .PluginNames}}	{{$pluginName | ToCamel | FirstUpper}}() {{$pluginName | ToCamel | FirstUpper}}
{{end}}
}

type plugins struct {
{{range $pluginName := .PluginNames}}	{{$pluginName | ToCamel | FirstLower}} {{$pluginName | ToCamel | FirstUpper}}
{{end}}
}

{{range $pluginName := .PluginNames}}func (p *plugins) {{$pluginName | ToCamel | FirstUpper}}() {{$pluginName | ToCamel | FirstUpper}} {
	return p.{{$pluginName | ToCamel | FirstLower}}
}

{{end}}

func NewPlugins(c restc.Client) Plugins {
	return &plugins{
{{range $pluginName := .PluginNames}}		{{$pluginName | ToCamel | FirstLower}}: New{{$pluginName | ToCamel | FirstUpper}}(c),
{{end}}
	}
}