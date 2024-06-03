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
}

type GatewayConf struct {
	gateway.GatewayConf
}

type LogConf struct {
	logx.LogConf
}