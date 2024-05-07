syntax = "v1"

type pathRequest {
    Name string `path:"name"`
}

type response {
    Message string
}

@server(
    prefix: /api/v1
    group: hello
)
service {{ .APP }} {
    @handler HelloPathHandler
    get /hello/:name (pathRequest) returns (response)
}
