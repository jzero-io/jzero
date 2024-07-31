package middleware

import (
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterZrpc(z *zrpc.RpcServer) {
	z.AddUnaryInterceptors(ServerValidationUnaryInterceptor)
}

func RegisterGateway(g *gateway.Server) {
    httpx.SetErrorHandler(ErrorHandler)
    g.Use(ResponseHandler)
}
