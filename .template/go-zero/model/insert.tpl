func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
    statement, args := sqlbuilder.NewInsertBuilder().
            InsertInto(m.table).
            Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
            Values({{.expressionValues}}).Build()
	{{if .withCache}}{{.keys}}
    return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        if session != nil {
            return session.ExecCtx(ctx, statement, args...)
        }
		return conn.ExecCtx(ctx, statement, args...)
	}, {{.keyValues}}){{else}}if session != nil {
       return session.ExecCtx(ctx, statement, args...)
	}
	return m.conn.ExecCtx(ctx, statement, args...){{end}}
}
