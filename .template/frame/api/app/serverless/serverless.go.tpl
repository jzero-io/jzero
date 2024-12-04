{{ if has "serverless" .Features }}package serverless

import (
	"path/filepath"

	"{{ .Module }}/server/config"
	"{{ .Module }}/server/handler"
	"{{ .Module }}/server/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// Serverless please replace coreSvcCtx any type to real core svcCtx
func Serverless(server *rest.Server, coreSvcCtx any) {
	var c config.Config

	if err := conf.Load(filepath.Join("plugins", "{{ .DirName }}", "etc", "etc.yaml"), &c); err != nil {
		panic(err)
	}
	config.C = c

	svcCtx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, svcCtx)
}{{end}}