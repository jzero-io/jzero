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
	GrpcMaxConns       int    `json:",default=10000"`
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
