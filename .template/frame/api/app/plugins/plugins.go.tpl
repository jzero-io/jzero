{{ if has "serverless_core" .Features }}// Code generated by jzero. DO NOT EDIT.
package plugins

import (
	"{{ .Module }}/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func LoadPlugins(server *rest.Server, svcCtx *svc.ServiceContext) {}{{ end }}