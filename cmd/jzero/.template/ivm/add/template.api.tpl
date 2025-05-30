syntax = "v1"

info (
    go_package: "{{ .Group }}"
)

{{range $v := .Handlers | uniq}}type {{$v.Name | FirstUpper}}Request {}

type {{$v.Name | FirstUpper}}Response {}

{{end}}

@server (
    prefix: /api/v1
    group:  {{ .Group }}
)
service {{ .Service }} {
    {{range $v := .Handlers | uniq}}@handler {{$v.Name}}Handler
    {{$v.Verb}} /{{ $.Group }}/{{$v.Name | FirstLower}} ({{$v.Name | FirstUpper}}Request) returns ({{$v.Name | FirstUpper}}Response)

    {{end}}
}

