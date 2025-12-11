{{ if .Serverless }}{{else}}syntax = "v1"

info (
    go_package: "version"
    WrapCodeMsg: true
)

type VersionRequest {}

type VersionResponse {
    Version string `json:"version"`
    GoVersion string `json:"goVersion"`
    Commit string `json:"commit"`
    Date string `json:"date"`
}

@server(
    group: version
)
service {{ .APP | ToCamel }} {
    @handler Version
    get /api/version (VersionRequest) returns (VersionResponse)
}{{end}}