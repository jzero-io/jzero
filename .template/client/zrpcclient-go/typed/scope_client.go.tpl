package {{.Scope}}

import (
	"github.com/zeromicro/go-zero/zrpc"
	{{ range $v := .Services}}
	"{{ $.Module }}/typed/{{$.Scope}}/{{ $v | FirstLower }}"
	{{ end }}
)

type Interface interface {
	{{ range $v := .Services }}
	{{ $v | FirstUpper }}() {{ $v | FirstLower }}.{{ $v | FirstUpper }}
	{{ end }}
}

type Client struct {
	client zrpc.Client
}

func New(c zrpc.Client) *Client {
	return &Client{client: c}
}

{{ range $v := .Services }}
func (x *Client) {{ $v | FirstUpper }}() {{ $v | FirstLower }}.{{ $v | FirstUpper }} {
	return {{ $v | FirstLower }}.New{{ $v | FirstUpper }}(x.client)
}
{{ end }}