package worktabd

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jaronnie/worktab/worktabd/internal/config"
	"github.com/jaronnie/worktab/worktabd/internal/rest"
	"github.com/jaronnie/worktab/worktabd/internal/server"
	"github.com/jaronnie/worktab/worktabd/internal/svc"
	"github.com/jaronnie/worktab/worktabd/worktabdpb"
)

func StartworktabdZrpcServer(configFile string) {
	var c config.Config
	conf.MustLoad(viper.ConfigFileUsed(), &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		worktabdpb.RegisterWorktabdServer(grpcServer, server.NewWorktabdServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func StartworktabdGatewayServer() {
	httpAddress := viper.GetString("ListenOnHTTP")
	grpcAddress := viper.GetString("ListenOn")
	sock := viper.GetString("ListenOnSocket")

	if httpAddress == "" {
		httpAddress = "0.0.0.0:8090"
	}
	if grpcAddress == "" {
		grpcAddress = "0.0.0.0:9603"
	}

	g := gin.New()

	// load gin handler
	g = rest.Router(g)

	s := &http.Server{
		Addr:    httpAddress,
		Handler: g,
	}
	fmt.Printf("Starting http server at %s...\n", httpAddress)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// unix socket server
	if sock != "" {
		os.Remove(sock)
		fmt.Printf("Starting unix socket server at %s...\n", sock)
		unixListener, err := net.Listen("unix", sock)
		if err != nil {
			panic(err)
		}
		go func() {
			if err := http.Serve(unixListener, g); err != nil {
				panic(err)
			}
		}()
	}
}
