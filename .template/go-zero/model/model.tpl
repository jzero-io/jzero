package {{.pkg}}
{{if .withCache}}
import (
    "context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/jzero-io/jzero-contrib/condition"
)
{{else}}

import (
    "context"

    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/jzero-io/jzero-contrib/condition"
)
{{end}}
var _ {{.upperStartCamelObject}}Model = (*custom{{.upperStartCamelObject}}Model)(nil)

type (
	// {{.upperStartCamelObject}}Model is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}Model interface {
		{{.lowerStartCamelObject}}Model
		{{if not .withCache}}WithSession(session sqlx.Session) {{.upperStartCamelObject}}Model{{end}}

	    BulkInsert(ctx context.Context, datas []*{{.upperStartCamelObject}}) error
        Find(ctx context.Context, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, error)
        Page(ctx context.Context, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, int64 ,error)
	}

	custom{{.upperStartCamelObject}}Model struct {
		*default{{.upperStartCamelObject}}Model
	}
)

// New{{.upperStartCamelObject}}Model returns a model for the database table.
func New{{.upperStartCamelObject}}Model(conn sqlx.SqlConn{{if .withCache}}, c cache.CacheConf, opts ...cache.Option{{end}}) {{.upperStartCamelObject}}Model {
	return &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(conn{{if .withCache}}, c, opts...{{end}}),
	}
}

{{if not .withCache}}
func (m *custom{{.upperStartCamelObject}}Model) WithSession(session sqlx.Session) {{.upperStartCamelObject}}Model {
    return New{{.upperStartCamelObject}}Model(sqlx.NewSqlConnFromSession(session))
}
{{end}}

