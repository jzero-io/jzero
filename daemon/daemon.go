package daemon

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/syncx"
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

	ctx := svc.NewServiceContext(c)
	start(ctx)
}

func start(ctx *svc.ServiceContext) {
	// print log to console if Log.Mode is file or volume
	middlewares.PrintLogToConsole(ctx.Config)

	// 保持与以往的版本兼容
	s := getZrpcServer(ctx.Config, ctx)

	middlewares.RateLimit = syncx.NewLimit(ctx.Config.Jzero.GrpcMaxConns)
	s.AddUnaryInterceptors(middlewares.GrpcRateLimitInterceptors)

	gw := gateway.MustNewServer(ctx.Config.Gateway)

	gw.Use(middlewares.WrapResponse)
	httpx.SetErrorHandler(middlewares.GrpcErrorHandler)

	// gw add routes
	handler.RegisterMyHandlers(gw.Server, ctx)

	// gw add api routes
	handler.RegisterHandlers(gw.Server, ctx)

	// listen unix
	var unixListener net.Listener
	var err error
	if ctx.Config.Jzero.ListenOnUnixSocket != "" {
		sock := ctx.Config.Jzero.ListenOnUnixSocket
		unixListener, err = net.Listen("unix", sock)
		if err != nil {
			panic(err)
		}
		go func() {
			fmt.Printf("Starting unix server at %s...\n", ctx.Config.Jzero.ListenOnUnixSocket)
			if err := http.Serve(unixListener, gw); err != nil {
				return
			}
		}()
	}

	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)

	go func() {
		fmt.Printf("Starting rpc server at %s...\n", ctx.Config.ListenOn)
		fmt.Printf("Starting gateway server at %s:%d...\n", ctx.Config.Gateway.Host, ctx.Config.Gateway.Port)
		group.Start()
	}()

	signalHandler(ctx, group, unixListener)
}

func signalHandler(ctx *svc.ServiceContext, serviceGroup *service.ServiceGroup, unixListener net.Listener) {
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Println("Waiting 1 second...\nStopping rpc server and gateway server")
			time.Sleep(time.Second)
			serviceGroup.Stop()
			if ctx.Config.Jzero.ListenOnUnixSocket != "" {
				fmt.Println("Stopping unix server")
				unixListener.Close()
				_ = os.Remove(ctx.Config.Jzero.ListenOnUnixSocket)
			}
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
