package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Zrpc    ZrpcConf
	Gateway GatewayConf

	JzeroConf
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

type JzeroConf struct {
	Log   LogConf
	Mysql MysqlConfig
}

type MysqlConfig struct {
	Username string `json:",optional"`
	Password string `json:",optional"`
	Address  string `json:",optional"`
	Database string `json:",optional"`
}
