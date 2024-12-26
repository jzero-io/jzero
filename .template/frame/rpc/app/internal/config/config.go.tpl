package config

import (
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Zrpc    ZrpcConf
	Log     LogConf

	Banner BannerConf
}

type ZrpcConf struct {
	zrpc.RpcServerConf
}

type LogConf struct {
	logx.LogConf
}

type BannerConf struct {
	Text     string `json:",default=JZERO"`
	Color    string `json:",default=green"`
	FontName string `json:",default=starwars,options=big|larry3d|starwars|standard"`
}
