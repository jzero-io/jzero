package config

import (
    "github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Zrpc zrpc.RpcServerConf
	Gateway gateway.GatewayConf

	{{ .APP | FirstUpper | ToCamel }}Conf
}

type {{ .APP | FirstUpper | ToCamel  }}Conf struct {
	GrpcMaxConns       int    `json:",default=10000"`

	Log logx.LogConf
    // only Log.Mode is file or volume take effect
	LogToConsole bool `json:",default=true"`
}
