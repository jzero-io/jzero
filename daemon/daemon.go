package daemon

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/jaronnie/jzero/daemon/internal/config"
	"github.com/jaronnie/jzero/daemon/internal/handler"
	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/jaronnie/jzero/daemon/middlewares"
)

func Start(cfgFile string) {
	var c config.Config
	conf.MustLoad(cfgFile, &c)
	go func() {
		start(c)
	}()
}

func start(c config.Config) {
	ctx := svc.NewServiceContext(c)
	s := getZrpcServer(c, ctx)

	gw := gateway.MustNewServer(c.Gateway)

	gw.Use(middlewares.WrapResponse)
	httpx.SetErrorHandler(middlewares.GrpcErrorHandler)

	// gw add routes
	handler.RegisterMyHandlers(gw.Server, ctx)

	// gw add api routes
	handler.RegisterHandlers(gw.Server, ctx)

	// listen unix
	if c.Jzero.ListenOnUnixSocket != "" {
		sock := c.Jzero.ListenOnUnixSocket
		_ = os.Remove(sock)
		unixListener, err := net.Listen("unix", sock)
		if err != nil {
			panic(err)
		}
		go func() {
			fmt.Printf("Starting unix server at %s...\n", c.Jzero.ListenOnUnixSocket)
			if err := http.Serve(unixListener, gw); err != nil {
				panic(err)
			}
		}()
	}

	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway server at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
	group.Start()
}
