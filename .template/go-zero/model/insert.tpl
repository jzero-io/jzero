func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
    ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        sql, args := sqlbuilder.NewInsertBuilder().
            InsertInto(m.table).
            Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
            Values({{.expressionValues}}).Build()
		return conn.ExecCtx(ctx, sql, args...)
	}, {{.keyValues}}){{else}}sql, args := sqlbuilder.NewInsertBuilder().
                                          InsertInto(m.table).
                                          Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
                                          Values({{.expressionValues}}).Build()
    ret,err:=m.conn.ExecCtx(ctx, sql, args...){{end}}
	return ret,err
}
