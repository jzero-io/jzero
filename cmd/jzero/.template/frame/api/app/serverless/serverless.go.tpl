{{ if .Serverless }}package serverless

import (
	"path/filepath"

    "github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"{{ .Module }}/internal/config"
    "{{ .Module }}/internal/handler"
    "{{ .Module }}/internal/svc"
)

type Serverless struct {
	SvcCtx        *svc.ServiceContext                                   // 服务上下文
	HandlerFunc   func(server *rest.Server, svcCtx *svc.ServiceContext) // 服务路由
}

// New serverless function
func New() *Serverless {
	cc := configcenter.MustNewConfigCenter[config.Config](configcenter.Config{
		Type: "yaml",
	}, subscriber.MustNewFsnotifySubscriber(filepath.Join("plugins", "{{ .DirName }}", "etc", "etc.yaml"), subscriber.WithUseEnv(true)))

	svcCtx := svc.NewServiceContext(cc)

	return &Serverless{
		SvcCtx:      svcCtx,
		HandlerFunc: handler.RegisterHandlers,
	}
}{{end}}