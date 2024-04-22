package daemon

import (
	"fmt"
	"net"
    "net/http"
    "os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/rest/httpx"

	"{{ .Module }}/daemon/internal/config"
	"{{ .Module }}/daemon/internal/handler"
	"{{ .Module }}/daemon/internal/svc"
	"{{ .Module }}/daemon/middlewares"
)

func Start(cfgFile string) {
	conf.MustLoad(cfgFile, &config.C)
	go func() {
	    // print log to console if Log.Mode is file or volume
	    middlewares.PrintLogToConsole(config.C)
		start()
	}()
}

func start() {
	ctx := svc.NewServiceContext(config.C)
	s := getZrpcServer(config.C, ctx)

	middlewares.RateLimit = syncx.NewLimit(config.C.{{ .APP | FirstUpper }}.GrpcMaxConns)
    s.AddUnaryInterceptors(middlewares.GrpcRateLimitInterceptors)

	gw := gateway.MustNewServer(config.C.Gateway)

	// add middlewares
	gw.Use(middlewares.WrapResponse)

	// error handler
	httpx.SetErrorHandler(middlewares.GrpcErrorHandler)

	// gw add your routes
	handler.RegisterMyHandlers(gw.Server, ctx)

	// gw add go-zero api framework routes
	handler.RegisterHandlers(gw.Server, ctx)

	// listen unix
	if config.C.{{ .APP | FirstUpper }}.ListenOnUnixSocket != "" {
	    sock := config.C.{{ .APP | FirstUpper }}.ListenOnUnixSocket
	    _ = os.Remove(sock)
	    unixListener, err := net.Listen("unix", sock)
	    if err != nil {
		    panic(err)
	    }
	    go func() {
	        fmt.Printf("Starting unix server at %s...\n", config.C.{{ .APP | FirstUpper }}.ListenOnUnixSocket)
		    if err := http.Serve(unixListener, gw); err != nil {
			    panic(err)
		    }
	    }()
	}

	group := service.NewServiceGroup()
	group.Add(s)
	group.Add(gw)

	fmt.Printf("Starting rpc server at %s...\n", config.C.ListenOn)
	fmt.Printf("Starting gateway server at %s:%d...\n", config.C.Gateway.Host, config.C.Gateway.Port)
	group.Start()
}
