package svc

import (
	"{{ .Module }}/internal/config"
	"{{ .Module }}/internal/custom"

	configurator "github.com/zeromicro/go-zero/core/configcenter"
)

type ServiceContext struct {
	Config config.Config

	Custom *custom.Custom
}

func NewServiceContext(c config.Config, cc configurator.Configurator[config.Config]) *ServiceContext {
    sc := &ServiceContext{
		Config: c,
		Custom: custom.New(),
	}
	sc.DynamicConfListener(cc)
	return sc
}
