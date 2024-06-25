syntax = "v1"

type List{{ .GroupCamel }}Request {}

type List{{ .GroupCamel }}Response {}

type Create{{ .GroupCamel }}Request {}

type Create{{ .GroupCamel }}Response {}

type Edit{{ .GroupCamel }}Request {}

type Edit{{ .GroupCamel }}Response {}

type Delete{{ .GroupCamel }}Request {}

type Delete{{ .GroupCamel }}Response {}

type Get{{ .GroupCamel }}Request {}

type Get{{ .GroupCamel }}Response {}

@server (
    prefix: /api/v1
    group:  {{ .Group }}
)
service {{ .Service }} {
    @handler ListHandler
    get /{{ .Group }}/list (List{{ .GroupCamel }}Request) returns (List{{ .GroupCamel }}Response)

    @handler CreateHandler
    get /{{ .Group }}/create (Create{{ .GroupCamel }}Request) returns (Create{{ .GroupCamel }}Response)

    @handler EditHandler
    get /{{ .Group }}/edit (Edit{{ .GroupCamel }}Request) returns (Edit{{ .GroupCamel }}Response)

    @handler DeleteHandler
    get /{{ .Group }}/delete (Delete{{ .GroupCamel }}Request) returns (Delete{{ .GroupCamel }}Response)

    @handler GetHandler
    get /{{ .Group }} (Get{{ .GroupCamel }}Request) returns (Get{{ .GroupCamel }}Response)
}

