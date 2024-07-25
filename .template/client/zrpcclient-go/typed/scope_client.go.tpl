package {{.Scope}}

import (
	"github.com/zeromicro/go-zero/zrpc"
	{{ range $v := .Services}}
	"{{ $.Module }}/typed/{{$.Scope}}/{{ $v | lower }}"
	{{ end }}
)

type Interface interface {
	{{ range $v := .Services }}
	{{ $v | FirstUpper | ToCamel }}() {{ $v | lower }}.{{ $v | FirstUpper | ToCamel }}
	{{ end }}
}

type Client struct {
	client zrpc.Client
}

func New(c zrpc.Client) *Client {
	return &Client{client: c}
}

{{ range $v := .Services }}
func (x *Client) {{ $v | FirstUpper | ToCamel }}() {{ $v | lower }}.{{ $v | FirstUpper | ToCamel }} {
	return {{ $v | lower }}.New{{ $v | FirstUpper | ToCamel }}(x.client)
}
{{ end }}