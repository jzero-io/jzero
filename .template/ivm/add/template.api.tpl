syntax = "v1"

{{ if .SplitApiTypesDir }}
info (
    go_package: "{{ .Group }}"
)
{{end}}

{{range $v := .Handlers | uniq}}type {{$v.Name | FirstUpper}}{{ if $.SplitApiTypesDir }}{{else}}{{ $.GroupCamel }}{{end}}Request {}

type {{$v.Name | FirstUpper}}{{ if $.SplitApiTypesDir }}{{else}}{{ $.GroupCamel }}{{end}}Response {}

{{end}}

@server (
    prefix: /api/v1
    group:  {{ .Group }}
)
service {{ .Service }} {
    {{range $v := .Handlers | uniq}}@handler {{$v.Name}}Handler
    {{$v.Verb}} /{{ $.Group }}/{{$v.Name | FirstLower}} ({{$v.Name | FirstUpper}}{{ if $.SplitApiTypesDir }}{{else}}{{ $.GroupCamel }}{{end}}Request) returns ({{$v.Name | FirstUpper}}{{ if $.SplitApiTypesDir }}{{else}}{{ $.GroupCamel }}{{end}}Response)

    {{end}}
}

