{{ if .Serverless }}package serverless

import (
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"google.golang.org/grpc"

	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/global"
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
	cc := configurator.MustNewConfigCenter[config.Config](configurator.Config{
		Type: "yaml",
	}, subscriber.MustNewFsnotifySubscriber(filepath.Join("plugins", "{{ .DirName }}", "etc", "etc.yaml"), subscriber.WithUseEnv(true)))

	svcCtx := svc.NewServiceContext(cc)
	global.ServiceContext = *svcCtx

	// get protoSets
	protoSets, err := pb.WriteToLocal(pb.Embed)
	logx.Must(err)

	return &Serverless{
		SvcCtx:             svcCtx,
		RegisterZrpcServer: server.RegisterZrpcServer,
		ProtoSets:          protoSets,
	}
}{{end}}