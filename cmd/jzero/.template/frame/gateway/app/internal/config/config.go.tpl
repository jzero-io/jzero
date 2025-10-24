package config

import (
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/gateway"
    "github.com/zeromicro/go-zero/zrpc"
    {{ if has "model" .Features }}"github.com/zeromicro/go-zero/core/stores/sqlx"{{ end }}
    {{ if has "redis" .Features }}"github.com/zeromicro/go-zero/core/stores/redis"{{ end }}
)

type Config struct {
	Zrpc    ZrpcConf
	Gateway GatewayConf
	Log     LogConf
	{{ if has "model" .Features }}Sqlx SqlxConf{{ end }}
    {{ if has "redis" .Features }}Redis RedisConf{{ end }}
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

{{ if has "model" .Features }}type SqlxConf struct {
	sqlx.SqlConf
}{{ end }}
{{ if has "redis" .Features }}type RedisConf struct {
    redis.RedisConf
}{{ end }}

type BannerConf struct {
	Text     string `json:",default=JZERO"`
	Color    string `json:",default=green"`
	FontName string `json:",default=starwars,options=big|larry3d|starwars|standard"`
}
