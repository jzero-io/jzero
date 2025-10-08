{{ if .Serverless }}package serverless

import (
	"path/filepath"

    "github.com/jzero-io/jzero/core/configcenter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"google.golang.org/grpc"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/desc/pb"
	"{{ .Module }}/internal/server"
	"{{ .Module }}/internal/svc"
)

type Serverless struct {
	SvcCtx             *svc.ServiceContext // 服务上下文
	RegisterZrpcServer func(grpcServer *grpc.Server, ctx *svc.ServiceContext)
	ProtoSets          []string
}

func New() *Serverless {
	cc := configcenter.MustNewConfigCenter[config.Config](configcenter.Config{
		Type: "yaml",
	}, subscriber.MustNewFsnotifySubscriber(filepath.Join("plugins", "{{ .DirName }}", "etc", "etc.yaml"), subscriber.WithUseEnv(true)))

	svcCtx := svc.NewServiceContext(cc)

	// get protoSets
	protoSets, err := pb.WriteToLocal(pb.Embed)
	logx.Must(err)

	return &Serverless{
		SvcCtx:             svcCtx,
		RegisterZrpcServer: server.RegisterZrpcServer,
		ProtoSets:          protoSets,
	}
}{{end}}