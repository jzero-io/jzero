package handler

// myroutes. 非框架生成的 routes. 建议能用框架自带解决的就用自带的功能!

import (
	"github.com/jaronnie/jzero/daemon/internal/svc"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func RegisterMyHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// register static embed handler
	registerStaticEmbedHandler(server, serverCtx)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/v1.0/health",
				Handler: HealthHandler(serverCtx),
			},
		},
	)
}
