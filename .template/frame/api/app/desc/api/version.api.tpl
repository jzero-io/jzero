syntax = "v1"

info (
	go_package: "version"
)

type GetRequest {}

type GetResponse {
    version string `json:"version"`
    goVersion string `json:"goVersion"`
    commit string `json:"commit"`
    date string `json:"date"`
}

@server(
    prefix: /api{{ if has "serverless" .Features }}/{{ .APP }}{{end}}/v1
    group: version
)
service {{ .APP | ToCamel }} {
    @handler Get
    get /version (GetRequest) returns (GetResponse)
}
