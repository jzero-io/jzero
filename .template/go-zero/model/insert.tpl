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

func (m *custom{{.upperStartCamelObject}}Model) createBuilder(build sqlbuilder.InsertBuilder) *sqlbuilder.InsertBuilder {
    return build.InsertInto(m.table)
}

func (m *custom{{.upperStartCamelObject}}Model) BulkInsert(ctx context.Context, datas []*{{.upperStartCamelObject}}) error {
    builder := sqlbuilder.NewInsertBuilder()
    builder.Cols({{.lowerStartCamelObject}}RowsExpectAutoSet)
    for _, data := range datas {
        builder.Values({{.expressionValues}})
    }
    sql, args := m.createBuilder(*builder).Build()
    sql = strings.ReplaceAll(sql, "`", "")
    _, err:= m.conn.ExecCtx(ctx, sql, args...)
    return err
}
