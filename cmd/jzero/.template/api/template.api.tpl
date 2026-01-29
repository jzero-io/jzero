syntax = "v1"

info (
    go_package: "{{ .Group }}"
)

type CreateRequest {}

type CreateResponse {}

@server (
    prefix: /api/{{ .Group }}
    group:  {{ .Group }}
    compact_handler: true
)
service {{ .Service }} {
    @handler Create
    post / (CreateRequest) returns (CreateResponse)
}