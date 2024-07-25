package {{ .APP }}

import (
	"github.com/zeromicro/go-zero/zrpc"

	{{range $v := .Scopes}}"{{$.Module}}/typed/{{$v | lower}}"{{end}}
)

type Opt func(client *Clientset)

{{ range $v := .Scopes}}
func With{{ $v | FirstUpper | ToCamel }}Client(cli zrpc.Client) Opt {
	return func(client *Clientset) {
		client.{{ $v | ToCamel }} = {{ $v }}.New(cli)
	}
}
{{ end}}

