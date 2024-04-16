package daemon

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jaronnie/jzero/daemon/internal/config"
	credentialsvr "github.com/jaronnie/jzero/daemon/internal/server/credential"
	credentialsvrv2 "github.com/jaronnie/jzero/daemon/internal/server/credentialv2"
	machinesvr "github.com/jaronnie/jzero/daemon/internal/server/machine"
	machinesvrv2 "github.com/jaronnie/jzero/daemon/internal/server/machinev2"
	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/jaronnie/jzero/daemon/pb/credentialpb"
	"github.com/jaronnie/jzero/daemon/pb/machinepb"
)

func getZrpcServer(c config.Config, ctx *svc.ServiceContext) *zrpc.RpcServer {
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		credentialpb.RegisterCredentialServer(grpcServer, credentialsvr.NewCredentialServer(ctx))
		credentialpb.RegisterCredentialv2Server(grpcServer, credentialsvrv2.NewCredentialv2Server(ctx))
		machinepb.RegisterMachineServer(grpcServer, machinesvr.NewMachineServer(ctx))
		machinepb.RegisterMachinev2Server(grpcServer, machinesvrv2.NewMachinev2Server(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	return s
}
