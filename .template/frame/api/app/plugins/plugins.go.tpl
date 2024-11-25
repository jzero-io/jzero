{{ if has "serverless_core" .Features }}package plugins

import (
	"github.com/zeromicro/go-zero/rest"
)

func LoadPlugins(server *rest.Server) {
	// server.AddRoutes(serverless.Routes())
}{{ end }}