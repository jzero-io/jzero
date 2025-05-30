syntax = "v1"

type GetRequest {}

type GetResponse {
    Version string `json:"version"`
    GoVersion string `json:"goVersion"`
    Commit string `json:"commit"`
    Date string `json:"date"`
}

@server(
    prefix: /api/v1{{ if has "serverless" .Features }}/{{ .APP }}{{end}}
)
service {{ .APP | ToCamel }} {
    @handler Get
    get /version (GetRequest) returns (GetResponse)
}
