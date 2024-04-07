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

	// go func() {
	// 	startworktabdRestServer(c)
	// }()
}

func startworktabdZrpcServer(c config.Config) {
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		worktabdpb.RegisterWorktabdServer(grpcServer, server.NewWorktabdServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	// s.Start()

	gw := gateway.MustNewServer(c.Gateway)
	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway server at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	group.Start()
}

// func startworktabdRestServer(c config.Config) {
// 	httpAddress := fmt.Sprintf("0.0.0.0:%s", viper.GetString("Gateway.Port"))

// 	if httpAddress == "" {
// 		httpAddress = "0.0.0.0:8090"
// 	}

// 	g := gin.New()
// 	// wrap grpc gateway handler
// 	handler := adapter.Wrap(func(h http.Handler) http.Handler {
// 		return gateway.MustNewServer(c.Gateway)
// 	})

// 	g.Use(handler)

// 	fmt.Printf("Starting rest server at %s...\n", httpAddress)
// 	go func() {
// 		if err := g.Run(fmt.Sprintf(":%d", c.Gateway.Port)); err != nil {
// 			panic(err)
// 		}
// 	}()
// }
