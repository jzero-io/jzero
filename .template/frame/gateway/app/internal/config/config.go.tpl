package config

import (
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/gateway"
    "github.com/zeromicro/go-zero/zrpc"
)

var C Config

type Config struct {
	Zrpc    ZrpcConf
	Gateway GatewayConf
	Log     LogConf

	Banner BannerConf
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

type BannerConf struct {
	Text     string `json:",default=JZERO"`
	Color    string `json:",default=green"`
	FontName string `json:",default=starwars,options=big|larry3d|starwars|standard"`
}
