package custom

type Custom struct{}

func New() *Custom {
	return &Custom{}
}

// Start Please add custom logic here.
func (c *Custom) Start() {}

// Stop Please add shut down logic here.
func (c *Custom) Stop() {}
