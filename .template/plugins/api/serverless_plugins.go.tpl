// Code generated by jzero. DO NOT EDIT.
package plugins

import (
	"github.com/zeromicro/go-zero/rest"

    "{{ .Module }}/internal/svc"
	{{range $v := .Plugins}}{{ $v.Path | base }} "{{ $v.Module }}/serverless"
	{{end}}
)

func LoadPlugins(server *rest.Server, svcCtx *svc.ServiceContext) {
	{{ range $v := .Plugins }}
	{
        serverless := {{ $v.Path | base }}.New(svcCtx)
        serverless.HandlerFunc(server, serverless.SvcCtx)
    }
	{{end}}
}