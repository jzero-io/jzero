package custom

import (
	"os"

	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"

	"{{.Module}}/internal/config"
	"{{.Module}}/internal/global"
)

type Custom struct {
	ZrpcServer *zrpc.RpcServer
	GatewayServer *gateway.Server
}

func New(zrpcServer *zrpc.RpcServer, gatewayServer *gateway.Server) *Custom {
	return &Custom{
		ZrpcServer: zrpcServer,
		GatewayServer: gatewayServer,
	}
}

func (c *Custom) Init() {
    c.AddRoutes(c.GatewayServer)
}

// Start Please add custom logic here.
func (c *Custom) Start() {}

// Stop Please add shut down logic here.
func (c *Custom) Stop() {}
