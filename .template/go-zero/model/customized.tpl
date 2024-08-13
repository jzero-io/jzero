func (m *custom{{.upperStartCamelObject}}Model) Find(ctx context.Context, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
	sb := sqlbuilder.Select({{.lowerStartCamelObject}}FieldNames...).From(m.table)
	condition.Apply(sb, conds...)
	sql, args := sb.Build()

	var resp []*{{.upperStartCamelObject}}
	err := m.conn.QueryRowsCtx(ctx, &resp, sql, args...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *custom{{.upperStartCamelObject}}Model) Page(ctx context.Context, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, int64 ,error) {
	sb := sqlbuilder.Select({{.lowerStartCamelObject}}FieldNames...).From(m.table)
	countsb := sqlbuilder.Select("count(*)").From(m.table)

	condition.Apply(sb, conds...)
	condition.Apply(countsb, conds...)

	var resp []*{{.upperStartCamelObject}}

	sql, args := sb.Build()
	err := m.conn.QueryRowsCtx(ctx, &resp, sql, args...)
	if err != nil {
		return nil, 0, err
	}

	// get total
    var total int64
    sql, args = countsb.Build()
    err = m.conn.QueryRowCtx(ctx, &total, sql, args...)
    if err != nil {
    	return nil, 0, err
    }

	return resp, total, nil
}