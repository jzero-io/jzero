package svc

import (
	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/custom"
	"{{ .Module }}/internal/middleware"
)

type ServiceContext struct {
	Config config.Config
	Custom *custom.Custom

	// middleware
	middleware.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Custom: custom.New(),
		Middleware: middleware.NewMiddleware(),
	}
}
