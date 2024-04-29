package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/jzero-io/jzero/daemon/internal/config"
)

type ServiceContext struct {
	Config  config.Config
	SqlConn sqlx.SqlConn
}

func NewServiceContext(c config.Config, sqlConn sqlx.SqlConn) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		SqlConn: sqlConn,
	}
}
