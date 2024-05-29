package config

import (
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/gateway"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Zrpc    ZrpcConf
	Gateway GatewayConf
	Log     LogConf
}

type ZrpcConf struct {
	zrpc.RpcServerConf

	MaxConns int `json:",default=10000"`
}

type GatewayConf struct {
	gateway.GatewayConf

	ListenOnUnixSocket string `json:",optional"`
}

type LogConf struct {
	logx.LogConf
	// only Log.Mode is file or volume take effect
	LogToConsole bool `json:",default=true"`
}