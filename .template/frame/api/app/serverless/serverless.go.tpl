{{ if has "serverless" .Features }}package serverless

import (
	"path/filepath"

	"{{ .Module }}/server/config"
	"{{ .Module }}/server/handler"
	"{{ .Module }}/server/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

type Serverless struct {
	SvcCtx        *svc.ServiceContext                                   // 服务上下文
	HandlerFunc   func(server *rest.Server, svcCtx *svc.ServiceContext) // 服务路由
}

// Serverless please replace coreSvcCtx any type to real core svcCtx
func New(coreSvcCtx any) *Serverless {
	var c config.Config

	if err := conf.Load(filepath.Join("plugins", "{{ .DirName }}", "etc", "etc.yaml"), &c); err != nil {
		panic(err)
	}
	config.C = c

	svcCtx := svc.NewServiceContext(c)
	return &Serverless{
    	SvcCtx:        svcCtx,
    	HandlerFunc:   handler.RegisterHandlers,
    }
}{{end}}