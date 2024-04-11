package jzerod

import (
	"fmt"
	"net"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jaronnie/jzero/jzerod/internal/config"
	"github.com/jaronnie/jzero/jzerod/internal/handler"
	credentialsvr "github.com/jaronnie/jzero/jzerod/internal/server/credential"
	credentialsvrv2 "github.com/jaronnie/jzero/jzerod/internal/server/credentialv2"
	machinesvr "github.com/jaronnie/jzero/jzerod/internal/server/machine"
	machinesvrv2 "github.com/jaronnie/jzero/jzerod/internal/server/machinev2"
	"github.com/jaronnie/jzero/jzerod/internal/svc"
	"github.com/jaronnie/jzero/jzerod/pb/credentialpb"
	"github.com/jaronnie/jzero/jzerod/pb/machinepb"
)

func StartJzeroDaemon(cfgFile string) {
	var c config.Config
	conf.MustLoad(cfgFile, &c)
	go func() {
		startJzerodZrpcServer(c)
	}()
}

func startJzerodZrpcServer(c config.Config) {
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		credentialpb.RegisterCredentialServer(grpcServer, credentialsvr.NewCredentialServer(ctx))
		credentialpb.RegisterCredentialv2Server(grpcServer, credentialsvrv2.NewCredentialv2Server(ctx))
		machinepb.RegisterMachineServer(grpcServer, machinesvr.NewMachineServer(ctx))
		machinepb.RegisterMachinev2Server(grpcServer, machinesvrv2.NewMachinev2Server(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	gw := gateway.MustNewServer(c.Gateway)
	// gw add routes
	handler.RegisterMyHandlers(gw.Server, ctx)

	// gw add api routes
	handler.RegisterHandlers(gw.Server, ctx)

	// listen unix
	sock := "./jzero.sock"
	unixListener, err := net.Listen("unix", sock)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := http.Serve(unixListener, gw); err != nil {
			panic(err)
		}
	}()

	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway server at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	group.Start()
}
