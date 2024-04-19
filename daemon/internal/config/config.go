package config

import (
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Gateway gateway.GatewayConf

	Jzero JzeroConfig
}

type JzeroConfig struct {
	ListenOnUnixSocket string `json:",optional"`
}
