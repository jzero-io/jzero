package modelx

import (
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	DatabaseTypeMysql    = "mysql"
	DatabaseTypePostgres = "postgres"
)

type ModelConf struct {
	DatabaseType string `json:"databaseType,default=mysql,options=mysql|sqlite|postgres"`
	DatabaseUrl  string `json:"databaseUrl,optional"`

	Mysql    MysqlConf    `json:"mysql,optional"`
	Postgres PostgresConf `json:"postgres,optional"`
}

type MysqlConf struct {
	Host     string `json:"host,default=localhost"`
	Port     int    `json:"port,default=3306"`
	Username string `json:"username,default=root"`
	Password string `json:"password,default=123456"`
	Database string `json:"database,default=jzero"`
}

type PostgresConf struct {
	Host     string `json:"host,default=localhost"`
	Port     int    `json:"port,default=5432"`
	Username string `json:"username,default=root"`
	Password string `json:"password,default=123456"`
	Database string `json:"database,default=jzero"`
	SslMode  string `json:"sslMode,default=disable"`
}

func BuildDataSource(c ModelConf) string {
	switch c.DatabaseType {
	case DatabaseTypeMysql:
		sqlbuilder.DefaultFlavor = sqlbuilder.MySQL
		if c.DatabaseUrl != "" {
			return c.DatabaseUrl
		}
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Mysql.Username,
			c.Mysql.Password,
			c.Mysql.Host+":"+cast.ToString(c.Mysql.Port),
			c.Mysql.Database)
	case DatabaseTypePostgres:
		sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
		if c.DatabaseUrl != "" {
			return c.DatabaseUrl
		}
		return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
			c.Postgres.Username,
			c.Postgres.Password,
			c.Postgres.Host+":"+cast.ToString(c.Postgres.Port),
			c.Postgres.Database,
			c.Postgres.SslMode,
		)
	}
	return ""
}

func MustSqlConn(c ModelConf) sqlx.SqlConn {
	var sqlConn sqlx.SqlConn

	switch c.DatabaseType {
	case DatabaseTypeMysql:
		sqlConn = sqlx.NewMysql(BuildDataSource(c))
	case DatabaseTypePostgres:
		sqlConn = postgres.New(BuildDataSource(c))
	default:
		panic(fmt.Sprintf("not supported database type: %s", c.DatabaseType))
	}

	db, err := sqlConn.RawDB()
	logx.Must(err)
	logx.Must(db.Ping())
	return sqlConn
}
