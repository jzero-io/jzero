package modelx

import (
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	_ "modernc.org/sqlite"

	"github.com/jzero-io/jzero/core/stores/cache"
)

func MustNewConn(c sqlx.SqlConf) sqlx.SqlConn {
	sqlConn := sqlx.MustNewConn(c)
	db, err := sqlConn.RawDB()
	logx.Must(err)
	err = db.Ping()
	logx.Must(err)

	sqlbuilder.DefaultFlavor = getSqlbuilderFlavor(c.DriverName)
	return sqlConn
}

func MustNewConnAndSqlbuilderFlavor(c sqlx.SqlConf) (sqlx.SqlConn, sqlbuilder.Flavor) {
	sqlConn := sqlx.MustNewConn(c)
	db, err := sqlConn.RawDB()
	logx.Must(err)
	err = db.Ping()
	logx.Must(err)

	return sqlConn, getSqlbuilderFlavor(c.DriverName)
}

// NewConnWithCache returns a CachedConn with a custom cache.
func NewConnWithCache(db sqlx.SqlConn, c cache.Cache) sqlc.CachedConn {
	return sqlc.NewConnWithCache(db, c)
}

func getSqlbuilderFlavor(driverName string) sqlbuilder.Flavor {
	switch driverName {
	case "mysql":
		return sqlbuilder.MySQL
	case "pgx":
		return sqlbuilder.PostgreSQL
	case "sqlite":
		return sqlbuilder.SQLite
	default:
		return sqlbuilder.DefaultFlavor
	}
}
