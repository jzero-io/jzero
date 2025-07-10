package custom

import (
	"github.com/jzero-io/jzero/core/swaggerv2"
	"github.com/zeromicro/go-zero/gateway"
)

func (c *Custom) AddRoutes(gatewayServer *gateway.Server) {
	// gatewayServer add swagger routes. If you do not want it, you can delete this line
	swaggerv2.RegisterRoutes(gatewayServer.Server)

	// add custom route
	// gatewayServer.AddRoute(rest.Route{})
}
