func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, newData.{{.upperStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    _, {{if .containsIndexCache}}err{{else}}err:{{end}}= m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        sb := sqlbuilder.Update(m.table)
        split := strings.Split({{.lowerStartCamelObject}}RowsExpectAutoSet, ",")
        var assigns []string
        for _, s := range split {
           assigns = append(assigns, sb.Assign(s, nil))
        }
        sb.Set(assigns...)
        sb.Where(sb.EQ("{{.originalPrimaryKey}}", nil))
        sql, _ := sb.Build()
		return conn.ExecCtx(ctx, sql, {{.expressionValues}})
	}, {{.keyValues}}){{else}} sb := sqlbuilder.Update(m.table)
	split := strings.Split({{.lowerStartCamelObject}}RowsExpectAutoSet, ",")
	var assigns []string
    for _, s := range split {
        assigns = append(assigns, sb.Assign(s, nil))
    }
    sb.Set(assigns...)
    sb.Where(sb.EQ("{{.originalPrimaryKey}}", nil))
    sql, _ := sb.Build()
    _,err:=m.conn.ExecCtx(ctx, sql, {{.expressionValues}}){{end}}
	return err
}
