package config

import (
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Gateway gateway.GatewayConf

	{{ .APP | FirstUpper }} {{ .APP | FirstUpper }}Config
}

type {{ .APP | FirstUpper }}Config struct {
	ListenOnUnixSocket string `json:",optional"`
	GrpcMaxConns       int    `json:",default=10000"`
    // only Log.Mode is file or volume take effect
	LogToConsole bool `json:",default=true"`
}
