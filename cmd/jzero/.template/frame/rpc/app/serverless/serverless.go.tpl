{{ if .Serverless }}package serverless

import (
	"path/filepath"

    "github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"google.golang.org/grpc"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/server"
	"{{ .Module }}/internal/svc"
)

type Serverless struct {
	SvcCtx             *svc.ServiceContext // 服务上下文
	RegisterZrpcServer func(grpcServer *grpc.Server, ctx *svc.ServiceContext)
}

func New() *Serverless {
	cc := configcenter.MustNewConfigCenter[config.Config](configcenter.Config{
		Type: "yaml",
	}, subscriber.MustNewFsnotifySubscriber(filepath.Join("plugins", "{{ .DirName }}", "etc", "etc.yaml"), subscriber.WithUseEnv(true)))

	svcCtx := svc.NewServiceContext(cc)

	return &Serverless{
		SvcCtx:             svcCtx,
		RegisterZrpcServer: server.RegisterZrpcServer,
	}
}{{end}}