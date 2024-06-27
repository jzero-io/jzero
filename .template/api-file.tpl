syntax = "v1"

{{range $v := .Handlers | uniq}}type {{$v | FirstUpper}}{{ $.GroupCamel }}Request {}

type {{$v | FirstUpper}}{{ $.GroupCamel }}Response {}

{{end}}

@server (
    prefix: /api/v1
    group:  {{ .Group }}
)
service {{ .Service }} {
    {{range $v := .Handlers | uniq}}@handler {{$v}}Handler
    post /{{ $.Group }}/{{$v | FirstLower}} ({{$v | FirstUpper}}{{ $.GroupCamel }}Request) returns ({{$v | FirstUpper}}{{ $.GroupCamel }}Response)

    {{end}}
}

