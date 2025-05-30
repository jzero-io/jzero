package swaggerv2

func WithSwaggerHost(swaggerHost string) Opts {
	return func(config *swaggerConfig) {
		config.SwaggerHost = swaggerHost
	}
}

func WithSwaggerPath(swaggerPath string) Opts {
	return func(config *swaggerConfig) {
		config.SwaggerPath = swaggerPath
	}
}
