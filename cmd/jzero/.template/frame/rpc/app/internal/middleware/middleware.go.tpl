package middleware

import (
	"github.com/zeromicro/go-zero/zrpc"
)

func Register(z *zrpc.RpcServer) {
	z.AddUnaryInterceptors(NewValidator().UnaryServerMiddleware())
}
