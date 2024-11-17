func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, session sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
	sb.Where(sb.EQ("{{.originalPrimaryKey}}", {{.lowerStartCamelPrimaryKey}}))
	sb.Limit(1)
	sql, args := sb.Build()
	var resp {{.upperStartCamelObject}}
    var err error
    if session != nil {
    	err = session.QueryRowCtx(ctx, &resp, sql, args...)
	} else {
		err = m.conn.QueryRowCtx(ctx, &resp, sql, args...)
	}
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) FindOneWithCache(ctx context.Context, session sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.cachedConn.QueryRowCtx(ctx, &resp, {{.cacheKeyVariable}}, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
	    sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
	    sb.Where(sb.EQ("{{.originalPrimaryKey}}", {{.lowerStartCamelPrimaryKey}}))
        sql, args := sb.Build()
        if session != nil {
		    return session.QueryRowCtx(ctx, v, sql, args...)
	    }
		return conn.QueryRowCtx(ctx, v, sql, args...)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}return m.FindOne(ctx, session, {{.lowerStartCamelPrimaryKey}}){{end}}
}
