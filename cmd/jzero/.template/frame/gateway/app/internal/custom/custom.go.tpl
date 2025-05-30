package custom

import (
	"os"

	configurator "github.com/zeromicro/go-zero/core/configcenter"

	"{{.Module}}/internal/config"
)

type Custom struct {
	Config configurator.Configurator[config.Config]
}

func New(config configurator.Configurator[config.Config]) *Custom {
	return &Custom{Config: config}
}

// Start Please add custom logic here.
func (c *Custom) Start() {}

// Stop Please add shut down logic here.
func (c *Custom) Stop() {
	conf, err := c.Config.GetConfig()
	if err == nil {
		// remove temp pb file
		if len(conf.Gateway.Upstreams) > 0 {
			for _, p := range conf.Gateway.Upstreams[0].ProtoSets {
				_ = os.Remove(p)
			}
		}
	}
}
