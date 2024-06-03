syntax = "v1"

type GetVersionRequest {}

type GetVersionResponse {
    version string
    goVersion string
    commit string
    date string
}

@server(
    prefix: /api/v1
    group: version
)
service {{ .APP | ToCamel }} {
    @handler GetVersionHandler
    get /version (GetVersionRequest) returns (GetVersionResponse)
}
