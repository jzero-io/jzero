package custom

import (
	"github.com/jzero-io/jzero/core/swaggerv2"
	"github.com/zeromicro/go-zero/gateway"
)

func (c *Custom) AddRoutes(gw *gateway.Server) {
	// gw add swagger routes. If you do not want it, you can delete this line
	swaggerv2.RegisterRoutes(gw.Server)

	// add custom route
	// gw.AddRoute(rest.Route{})
}
