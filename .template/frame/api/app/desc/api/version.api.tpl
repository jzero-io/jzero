syntax = "v1"

type GetRequest {}

type GetResponse {
    version string `json:"version"`
    goVersion string `json:"goVersion"`
    commit string `json:"commit"`
    date string `json:"date"`
}

@server(
    prefix: /api/v1/{{ if has "serverless" .Features }}/{{ .APP }}{{end}}
    group: version
)
service {{ .APP | ToCamel }} {
    @handler Get
    get /version (GetRequest) returns (GetResponse)
}
