syntax = "v1"

info (
    go_package: "{{ .Group }}"
)

type CreateRequest {}

type CreateResponse {}

@server (
    prefix: /api
    group:  {{ .Group }}
)
service {{ .Service }} {
    @handler Create
    post / (CreateRequest) returns (CreateResponse)
}