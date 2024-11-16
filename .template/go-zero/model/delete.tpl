func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, session sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		sb := sqlbuilder.DeleteFrom(m.table)
        sb.Where(sb.EQ("{{.originalPrimaryKey}}", {{.lowerStartCamelPrimaryKey}}))
        sql, args := sb.Build()
        if session != nil {
			return session.ExecCtx(ctx, sql, args...)
		}
		return conn.ExecCtx(ctx, sql, args...)
	}, {{.keyValues}}){{else}}sb := sqlbuilder.DeleteFrom(m.table)
		 sb.Where(sb.EQ("{{.originalPrimaryKey}}", {{.lowerStartCamelPrimaryKey}}))
         sql, args := sb.Build()
         var err error
         if session != nil {
            _, err = session.ExecCtx(ctx, sql, args...)
		 }else{
			_, err = m.conn.ExecCtx(ctx, sql, args...)
		 }
         {{end}}
	return err
}
