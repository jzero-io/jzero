{{if .withCache}}
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, session sqlx.Session, {{.in}}) (*{{.upperStartCamelObject}}, error) {
	{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.cachedConn.QueryRowIndexCtx(ctx, &resp, {{.cacheKeyVariable}}, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
		condition.SelectByWhereRawSqlWithFlavor(m.flavor, sb, "{{.originalField}}", {{.lowerStartCamelField}})
		sb.Limit(1)
        sql, args := sb.BuildWithFlavor(m.flavor)
        var err error

        if session != nil {
            err = session.QueryRowCtx(ctx, &resp, sql, args...)
        } else {
            err = conn.QueryRowCtx(ctx, &resp, sql, args...)
        }
		if err != nil {
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
}

{{if .withCacheEnabled}}
func (m *default{{.upperStartCamelObject}}Model) FindOneNoCacheBy{{.upperField}}(ctx context.Context, session sqlx.Session, {{.in}}) (*{{.upperStartCamelObject}}, error) {
	var resp {{.upperStartCamelObject}}
    var err error

	sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
	condition.SelectByWhereRawSql(sb, "{{.originalField}}", {{.lowerStartCamelField}})
    sb.Limit(1)

    sql, args := sb.BuildWithFlavor(m.flavor)

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
}{{end}}
{{else}}
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, session sqlx.Session, {{.in}}) (*{{.upperStartCamelObject}}, error) {
	var resp {{.upperStartCamelObject}}
    var err error

	sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
	condition.SelectByWhereRawSql(sb, "{{.originalField}}", {{.lowerStartCamelField}})
    sb.Limit(1)

    sql, args := sb.BuildWithFlavor(m.flavor)

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
{{end}}
