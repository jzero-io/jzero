func (m *custom{{.upperStartCamelObject}}Model) BulkInsert(ctx context.Context, datas []*{{.upperStartCamelObject}}) error {
    sb := sqlbuilder.InsertInto(m.table)
    sb.Cols({{.lowerStartCamelObject}}RowsExpectAutoSet)
    for _, data := range datas {
        sb.Values({{.expressionValues}})
    }
    sql, args := sb.Build()
    _, err:= m.conn.ExecCtx(ctx, sql, args...)
    return err
}

func (m *custom{{.upperStartCamelObject}}Model) FindByCondition(ctx context.Context, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
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

func (m *custom{{.upperStartCamelObject}}Model) FindOneByCondition(ctx context.Context, conds ...condition.Condition) (*{{.upperStartCamelObject}}, error) {
	sb := sqlbuilder.Select({{.lowerStartCamelObject}}FieldNames...).From(m.table)
	condition.Apply(sb, conds...)
	sb.Limit(1)
	sql, args := sb.Build()

	var resp {{.upperStartCamelObject}}
	err := m.conn.QueryRowCtx(ctx, &resp, sql, args...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *custom{{.upperStartCamelObject}}Model) PageByCondition(ctx context.Context, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, int64 ,error) {
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