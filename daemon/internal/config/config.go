package config

import (
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
)

var C Config

type Config struct {
	zrpc.RpcServerConf
	Gateway gateway.GatewayConf

	Jzero JzeroConfig
}

type JzeroConfig struct {
	ListenOnUnixSocket string `json:",optional"`
	GrpcMaxConns       int    `json:",default=10000"`
	// only Log.Mode is file or volume take effect
	LogToConsole bool `json:",default=true"`
}
