func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, session sqlx.Session, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error {
	sb := sqlbuilder.Update(m.table)
	split := strings.Split({{.lowerStartCamelObject}}RowsExpectAutoSet, ",")
	var assigns []string
    for _, s := range split {
        assigns = append(assigns, sb.Assign(s, nil))
    }
    sb.Set(assigns...)
    sb.Where(sb.EQ("{{.originalPrimaryKey}}", nil))
    statement, _ := sb.Build()

    var err error
    if session != nil{
        _, err = session.ExecCtx(ctx, statement, {{.expressionValues}})
    }else{
        _, err = m.conn.ExecCtx(ctx, statement, {{.expressionValues}})
    }
	return err
}

func (m *default{{.upperStartCamelObject}}Model) UpdateWithCache(ctx context.Context, session sqlx.Session, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err := m.FindOne(ctx, session, newData.{{.upperStartCamelPrimaryKey}})
	if err != nil{
		return err
	}
    {{end}}{{.keys}}
    _, {{if .containsIndexCache}}err{{else}}err :{{end}}= m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        sb := sqlbuilder.Update(m.table)
        split := strings.Split({{.lowerStartCamelObject}}RowsExpectAutoSet, ",")
        var assigns []string
        for _, s := range split {
           assigns = append(assigns, sb.Assign(s, nil))
        }
        sb.Set(assigns...)
        sb.Where(sb.EQ("{{.originalPrimaryKey}}", nil))
        statement, _ := sb.Build()
        if session != nil{
            return session.ExecCtx(ctx, statement, {{.expressionValues}})
        }
		return conn.ExecCtx(ctx, statement, {{.expressionValues}})
	}, {{.keyValues}})
	return err{{else}}return m.Update(ctx, session, {{if .containsIndexCache}}newData{{else}}data{{end}}){{end}}
}
