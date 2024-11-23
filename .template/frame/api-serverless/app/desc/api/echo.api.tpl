syntax = "v1"

info (
	go_package: "echo"
)

type EchoRequest {
    message string `form:"message"`
}

type EchoResponse {
    message string `json:"message"`
}

@server(
    prefix: /api/{{ .APP }}
    group: echo
)
service {{ .APP | ToCamel }} {
    @handler Echo
    get /echo (EchoRequest) returns (EchoResponse)
}
