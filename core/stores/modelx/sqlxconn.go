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

	setSqlbuilderFlavor(c.DriverName)
	return sqlConn
}

// NewConnWithCache returns a CachedConn with a custom cache.
func NewConnWithCache(db sqlx.SqlConn, c cache.Cache) sqlc.CachedConn {
	return sqlc.NewConnWithCache(db, c)
}

func setSqlbuilderFlavor(driverName string) {
	switch driverName {
	case "mysql":
		sqlbuilder.DefaultFlavor = sqlbuilder.MySQL
	case "pgx":
		sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	case "sqlite":
		sqlbuilder.DefaultFlavor = sqlbuilder.SQLite
	}
}
