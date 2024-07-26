package custom

import (
    "os"

    "{{.Module}}/internal/config"
)

type Custom struct{}

func New() *Custom {
	return &Custom{}
}

// Start Please add custom logic here.
func (c *Custom) Start() {}

// Stop Please add shut down logic here.
func (c *Custom) Stop() {
    // remove temp pb file
	if len(config.C.Gateway.Upstreams) > 0 {
		for _, p := range config.C.Gateway.Upstreams[0].ProtoSets {
			_ = os.Remove(p)
		}
	}
}
