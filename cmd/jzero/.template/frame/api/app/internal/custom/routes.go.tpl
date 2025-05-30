package custom

import (
	"github.com/jzero-io/jzero-contrib/swaggerv2"
	"github.com/zeromicro/go-zero/rest"
)

func (c *Custom) AddRoutes(server *rest.Server) {
	// server add swagger routes. If you do not want it, you can delete this line
	swaggerv2.RegisterRoutes(server)

	// add custom route
	// server.AddRoute(rest.Route{})
}
