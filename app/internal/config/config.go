package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Zrpc    zrpc.RpcServerConf
	Gateway gateway.GatewayConf

	JzeroConf
}

type JzeroConf struct {
	ListenOnUnixSocket string `json:",optional"`
	GrpcMaxConns       int    `json:",default=10000"`

	Log logx.LogConf
	// only Log.Mode is file or volume take effect
	LogToConsole bool `json:",default=true"`

	Mysql MysqlConfig
}

type MysqlConfig struct {
	Username string `json:",optional"`
	Password string `json:",optional"`
	Address  string `json:",optional"`
	Database string `json:",optional"`
}
