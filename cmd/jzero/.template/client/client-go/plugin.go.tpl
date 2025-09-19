{{define "methodDefine"}}{{.MethodName | FirstUpper}}({{if or .IsStreamServer .IsStreamClient }}{{else}}ctx context.Context,{{end}}{{ .Request.FullName }}) ({{if or .IsStreamServer .IsStreamClient}}*restc.Request{{else}}{{.Response.FullName}}{{end}}, error){{end}}
package plugins

import (
	{{if .UngroupedAPIs}}"context"

	{{end}}"github.com/jzero-io/jzero/core/restc"
{{range $resource := .Resources}}
	{{$.PluginName | ToCamel | FirstLower}}{{$resource | ToCamel | FirstUpper}} "{{$.Module}}/plugins/{{$.PluginName}}/typed/{{$resource | lower}}"{{end}}
{{range $v := .GoImportPaths | uniq}}"{{$v}}"
{{end}}
)

type {{.PluginName | ToCamel | FirstUpper}} interface {
{{range $resource := .Resources}}	{{$resource | ToCamel | FirstUpper}}() {{$.PluginName | ToCamel | FirstLower}}{{$resource | ToCamel | FirstUpper}}.{{$resource | ToCamel | FirstUpper}}
{{end}}{{range $k, $v := .UngroupedAPIs}}	// {{$v.MethodName | FirstUpper}} {{.Comments}}
	// {{$v.Method}}:{{$v.URL}}
	{{template "methodDefine" $v}}
{{end}}
}

type {{.PluginName | ToCamel | FirstLower}}Client struct {
	client restc.Client
{{range $resource := .Resources}}	{{$resource | ToCamel | FirstLower}} {{$.PluginName | ToCamel | FirstLower}}{{$resource | ToCamel | FirstUpper}}.{{$resource | ToCamel | FirstUpper}}
{{end}}
}

{{range $resource := .Resources}}func (x *{{$.PluginName | ToCamel | FirstLower}}Client) {{$resource | ToCamel | FirstUpper}}() {{$.PluginName | ToCamel | FirstLower}}{{$resource | ToCamel | FirstUpper}}.{{$resource | ToCamel | FirstUpper}} {
	return x.{{$resource | ToCamel | FirstLower}}
}

{{end}}{{range $k, $v := .UngroupedAPIs}}func (x *{{$.PluginName | ToCamel | FirstLower}}Client) {{$v.MethodName | FirstUpper}}({{if or $v.IsStreamServer $v.IsStreamClient }}{{else}}ctx context.Context,{{end}}{{ $v.Request.FullName }}) ({{if or $v.IsStreamServer $v.IsStreamClient}}*restc.Request{{else}}{{$v.Response.FullName}}{{end}}, error) {
	{{if or $v.IsStreamServer $v.IsStreamClient}}request := x.client.Verb("{{$v.Method}}").
		Path(
			"{{$v.URL}}",{{range $p := $v.PathParams}}
			restc.PathParam{Name: "{{$p.Name}}", Value: in.{{$p.GoName}}},{{end}}
		)
	return request, nil{{else}}var resp {{$v.Response.FullName}}
		err := x.client.Verb("{{$v.Method}}").
		Path(
			"{{$v.URL}}",{{range $p := $v.PathParams}}
			restc.PathParam{Name: "{{$p.Name}}", Value: in.{{$p.GoName}}},{{end}}
		).
		Params({{if eq $v.Request.Body "*"}}{{else}}{{range $q := $v.QueryParams}}
			restc.QueryParam{Name: "{{$q.Name}}", Value: in.{{$q.GoName}}},{{end}}{{end}}
		).
		{{ if or (eq $v.Method "GET") (eq $v.Method "DELETE") }}{{else}}Body({{if eq $v.Request.Body ""}}nil{{else if eq $v.Request.Body "*"}}in{{else if or (ne $v.Method "GET") (ne $v.Method "DELETE")}}in.{{$v.Request.RealBodyName}}{{else}}nil{{end}}).{{end}}
		Do(ctx).
		Into(&resp, {{if $v.WrapCodeMsg}}&restc.IntoOptions{
			WrapCodeMsg:        true,
			{{if $v.WrapCodeMsgMapping}}WrapCodeMsgMapping: restc.WrapCodeMsgMapping{
				Code: "{{$v.WrapCodeMsgMapping.Code}}",
				Data: "{{$v.WrapCodeMsgMapping.Data}}",
				Msg:  "{{$v.WrapCodeMsgMapping.Msg}}",
			},{{end}}
		}{{else}}nil{{end}})

	if err != nil {
		return nil, err
	}

	return resp, nil{{end}}
}
{{end}}

func New{{.PluginName | ToCamel | FirstUpper}}(c restc.Client) {{.PluginName | ToCamel | FirstUpper}} {
	return &{{.PluginName | ToCamel | FirstLower}}Client{
		client: c,
{{range $resource := .Resources}}		{{$resource | ToCamel | FirstLower}}: {{$.PluginName | ToCamel | FirstLower}}{{$resource | ToCamel | FirstUpper}}.New{{$resource | ToCamel | FirstUpper}}(c),
{{end}}
	}
}