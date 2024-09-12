syntax = "v1"

{{range $v := .Handlers | uniq}}type {{$v.Name | FirstUpper}}{{ $.GroupCamel }}Request {}

type {{$v.Name | FirstUpper}}{{ $.GroupCamel }}Response {}

{{end}}

@server (
    prefix: /api/v1
    group:  {{ .Group }}
)
service {{ .Service }} {
    {{range $v := .Handlers | uniq}}@handler {{$v.Name}}Handler
    {{$v.Verb}} /{{ $.Group }}/{{$v.Name | FirstLower}} ({{$v.Name | FirstUpper}}{{ $.GroupCamel }}Request) returns ({{$v.Name | FirstUpper}}{{ $.GroupCamel }}Response)

    {{end}}
}

