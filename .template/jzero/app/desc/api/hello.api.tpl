syntax = "v1"

type pathRequest {
    Name string `path:"name"`
}

type response {
    Message string `json:"message"`
}

@server(
    prefix: /api/v1
    group: hello
)
service {{ .APP | ToCamel }} {
    @handler HelloPathHandler
    get /hello/:name (pathRequest) returns (response)
}
