package custom

import (
	"github.com/zeromicro/go-zero/rest"
)

type Custom struct {
	Server *rest.Server
}

func New(server *rest.Server) *Custom {
	return &Custom{
		Server: server,
	}
}

// Init Please add custom logic here.
func (c *Custom) Init() {
	c.AddRoutes(c.Server)
}

// Start Please add custom logic here.
func (c *Custom) Start() {}

// Stop Please add shut down logic here.
func (c *Custom) Stop() {}
