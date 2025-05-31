package {{ .Package }}

import (
	"github.com/zeromicro/go-zero/zrpc"

	{{range $v := .Scopes}}"{{$.Module}}/{{ $.ClientDir }}"{{end}}
)

type Opt func(client *Clientset)

{{ range $v := .Scopes}}
func With{{ $v | FirstUpper | ToCamel }}Client(cli zrpc.Client) Opt {
	return func(client *Clientset) {
		client.{{ $v | ToCamel }} = {{ $v }}.New(cli)
	}
}
{{ end}}

