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

	"{{ .Module }}/daemon/internal/config"
	"{{ .Module }}/daemon/internal/handler"
	"{{ .Module }}/daemon/internal/svc"
	"{{ .Module }}/daemon/middlewares"
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

	// add middlewares
	gw.Use(middlewares.WrapResponse)

	// error handler
	httpx.SetErrorHandler(middlewares.GrpcErrorHandler)

	// gw add your routes
	handler.RegisterMyHandlers(gw.Server, ctx)

	// gw add go-zero api framework routes
	handler.RegisterHandlers(gw.Server, ctx)

	// listen unix
	if c.{{ .APP | FirstUpper }}.ListenOnUnixSocket != "" {
	    sock := c.{{ .APP | FirstUpper }}.ListenOnUnixSocket
	    _ = os.Remove(sock)
	    unixListener, err := net.Listen("unix", sock)
	    if err != nil {
		    panic(err)
	    }
	    go func() {
	        fmt.Printf("Starting unix server at %s...\n", c.{{ .APP | FirstUpper }}.ListenOnUnixSocket)
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
