{{ if .Serverless }}package serverless

import (
	"path/filepath"

	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/global"
    "{{ .Module }}/internal/handler"
    "{{ .Module }}/internal/svc"
)

type Serverless struct {
	SvcCtx        *svc.ServiceContext                                   // 服务上下文
	HandlerFunc   func(server *rest.Server, svcCtx *svc.ServiceContext) // 服务路由
}

// New serverless function
func New() *Serverless {
	cc := configurator.MustNewConfigCenter[config.Config](configurator.Config{
		Type: "yaml",
	}, subscriber.MustNewFsnotifySubscriber(filepath.Join("plugins", "{{ .DirName }}", "etc", "etc.yaml"), subscriber.WithUseEnv(true)))

	svcCtx := svc.NewServiceContext(cc)
	global.ServiceContext = *svcCtx

	return &Serverless{
		SvcCtx:      svcCtx,
		HandlerFunc: handler.RegisterHandlers,
	}
}{{end}}