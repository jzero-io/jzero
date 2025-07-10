syntax = "v1"

info (
    go_package: "version"
    version: "v1"
)

type VersionRequest {}

type VersionResponse {
    Version string `json:"version"`
    GoVersion string `json:"goVersion"`
    Commit string `json:"commit"`
    Date string `json:"date"`
}

@server(
    prefix: /api/v1{{ if .Serverless }}/{{ .APP }}{{end}}
    group: version
)
service {{ .APP | ToCamel }} {
    @handler Version
    get /version (VersionRequest) returns (VersionResponse)
}
