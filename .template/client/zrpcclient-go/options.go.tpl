package {{ .APP }}

import (
	"github.com/zeromicro/go-zero/zrpc"

	{{range $v := .Scopes}}"{{$.Module}}/typed/{{$v}}"{{end}}
)

type Opt func(client *Clientset)

{{ range $v := .Scopes}}
func With{{ $v | FirstUpper }}Client(cli zrpc.Client) Opt {
	return func(client *Clientset) {
		client.{{ $v }} = {{ $v }}.New(cli)
	}
}
{{ end}}

