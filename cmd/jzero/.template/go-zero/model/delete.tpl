{{if .withCache}}
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, session sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err := m.FindOne(ctx, session, {{.lowerStartCamelPrimaryKey}})
    	if err != nil{
    		return err
    	}

    {{end}}	{{.keys}}
        _, err {{if .containsIndexCache}}={{else}}:={{end}} m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
    		sb := sqlbuilder.DeleteFrom(m.table)
            sb.Where(sb.EQ(condition.QuoteWithFlavor(m.flavor, "{{.originalPrimaryKey}}"), {{.lowerStartCamelPrimaryKey}}))
            statement, args := sb.BuildWithFlavor(m.flavor)
            if session != nil {
    			return session.ExecCtx(ctx, statement, args...)
    		}
    		return conn.ExecCtx(ctx, statement, args...)
    	}, {{.keyValues}})
    	return err{{else}}return m.Delete(ctx, session, {{.lowerStartCamelPrimaryKey}}){{end}}
}
{{else}}
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, session sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	sb := sqlbuilder.DeleteFrom(m.table)
    sb.Where(sb.EQ(condition.QuoteWithFlavor(m.flavor, "{{.originalPrimaryKey}}"), {{.lowerStartCamelPrimaryKey}}))
    statement, args := sb.BuildWithFlavor(m.flavor)
    var err error
    if session != nil {
        _, err = session.ExecCtx(ctx, statement, args...)
    } else {
        _, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	return err
}
{{end}}