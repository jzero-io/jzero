func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
    statement, args := sqlbuilder.NewInsertBuilder().
            InsertInto(m.table).
            Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
            Values({{.expressionValues}}).Build()
	if session != nil {
       return session.ExecCtx(ctx, statement, args...)
	}
	return m.conn.ExecCtx(ctx, statement, args...)
}

func (m *default{{.upperStartCamelObject}}Model) InsertWithCache(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
	statement, args := sqlbuilder.NewInsertBuilder().
                InsertInto(m.table).
                Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
                Values({{.expressionValues}}).Build()
    return m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        if session != nil {
            return session.ExecCtx(ctx, statement, args...)
        }
		return conn.ExecCtx(ctx, statement, args...)
	}, {{.keyValues}}){{else}}return m.Insert(ctx, session, data){{end}}
}