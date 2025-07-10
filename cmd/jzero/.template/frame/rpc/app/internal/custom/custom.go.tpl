package custom

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Custom struct{
    ZrpcServer *zrpc.RpcServer
}

func New(zrpcServer *zrpc.RpcServer) *Custom {
	return &Custom{
		ZrpcServer: zrpcServer,
	}
}

func (c *Custom) Init() {}

// Start Please add custom logic here.
func (c *Custom) Start() {}

// Stop Please add shut down logic here.
func (c *Custom) Stop() {}
