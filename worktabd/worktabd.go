package worktabd

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jaronnie/worktab/worktabd/internal/config"
	"github.com/jaronnie/worktab/worktabd/internal/handler"
	"github.com/jaronnie/worktab/worktabd/internal/server"
	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/worktabdpb"
)

func StartWorktabDaemon(cfgFile string) {
	var c config.Config
	conf.MustLoad(cfgFile, &c)
	go func() {
		startworktabdZrpcServer(c)
	}()
}

func startworktabdZrpcServer(c config.Config) {
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		worktabdpb.RegisterWorktabdServer(grpcServer, server.NewWorktabdServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	gw := gateway.MustNewServer(c.Gateway)
	// gw add routes
	handler.RegisterMyHandlers(gw.Server, ctx)

	// gw add api routes
	handler.RegisterHandlers(gw.Server, ctx)

	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway server at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	group.Start()
}
