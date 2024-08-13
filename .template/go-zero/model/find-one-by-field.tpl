func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRowIndexCtx(ctx, &resp, {{.cacheKeyVariable}}, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
        // patch
		sb.Where(sb.EQ(strings.Split(strings.TrimSpace("{{.originalField}}"), "=")[0], {{.lowerStartCamelField}}))
		sb.Limit(1)
        sql, args := sb.Build()
		if err := conn.QueryRowCtx(ctx, &resp, sql, args...); err != nil {
			return nil, err
		}
		return resp.{{.upperStartCamelPrimaryKey}}, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{else}}var resp {{.upperStartCamelObject}}
	sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
    sb.Where(sb.EQ("{{.originalField}}", {{.lowerStartCamelField}}))
    sb.Limit(1)
    sql, args := sb.Build()
    err := m.conn.QueryRowCtx(ctx, &resp, sql, args...)

	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{end}}
