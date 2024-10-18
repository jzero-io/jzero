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
    prefix: /api/v1
    group: version
)
service {{ .APP | ToCamel }} {
    @handler GetHandler
    get /version (GetRequest) returns (GetResponse)
}
