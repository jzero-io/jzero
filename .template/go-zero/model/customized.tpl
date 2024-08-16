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
	condition.ApplySelect(sb, conds...)
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
	condition.ApplySelect(sb, conds...)
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

	condition.ApplySelect(sb, conds...)
	condition.ApplySelect(countsb, conds...)

	var resp []*{{.upperStartCamelObject}}

	sql, args := sb.Build()
	err := m.conn.QueryRowsCtx(ctx, &resp, sql, args...)
	if err != nil {
		return nil, 0, err
	}

    var total int64
    sql, args = countsb.Build()
    err = m.conn.QueryRowCtx(ctx, &total, sql, args...)
    if err != nil {
    	return nil, 0, err
    }

	return resp, total, nil
}

func (m *custom{{.upperStartCamelObject}}Model) UpdateFieldsByCondition(ctx context.Context, field map[string]any, conds ...condition.Condition) error {
    if field == nil {
        return nil
    }

	sb := sqlbuilder.Update(m.table)
	condition.ApplyUpdate(sb, conds...)

	var assigns []string
    for key, value := range field {
        assigns = append(assigns, sb.Assign(key, value))
    }
    sb.Set(assigns...)

	sql, args := sb.Build()
	_, err := m.conn.ExecCtx(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *custom{{.upperStartCamelObject}}Model) BulkDelete(ctx context.Context, conds ...condition.Condition) error {
    if len(conds) == 0 {
		return nil
	}
	sb := sqlbuilder.DeleteFrom(m.table)
	condition.ApplyDelete(sb, conds...)
	sql, args := sb.Build()
	_, err := m.conn.ExecCtx(ctx, sql, args...)
	return err
}