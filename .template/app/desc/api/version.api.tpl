syntax = "v1"

info (
	go_package: "version"
)

type GetVersionRequest {}

type GetVersionResponse {
    version string `json:"version"`
    goVersion string `json:"goVersion"`
    commit string `json:"commit"`
    date string `json:"date"`
}

@server(
    prefix: /api/v1
    group: version
)
service {{ .APP | ToCamel }} {
    @handler GetVersionHandler
    get /version (GetVersionRequest) returns (GetVersionResponse)
}