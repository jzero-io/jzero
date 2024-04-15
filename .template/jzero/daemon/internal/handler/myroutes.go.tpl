package handler

import (
	"{{ .Module }}/daemon/internal/svc"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

func RegisterMyHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
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
