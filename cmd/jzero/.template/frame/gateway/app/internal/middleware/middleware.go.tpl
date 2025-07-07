package middleware

import (
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func Register(z *zrpc.RpcServer, gw *gateway.Server) {
	z.AddUnaryInterceptors(NewValidator().UnaryServerMiddleware())
	z.AddUnaryInterceptors(WithValueMiddleware)

    httpx.SetErrorHandlerCtx(ErrorMiddleware)
    gw.Use(ResponseMiddleware)
}
