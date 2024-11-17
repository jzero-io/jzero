func (m *custom{{.upperStartCamelObject}}Model) BulkInsert(ctx context.Context, session sqlx.Session, datas []*{{.upperStartCamelObject}}) error {
    sb := sqlbuilder.InsertInto(m.table)
    sb.Cols({{.lowerStartCamelObject}}RowsExpectAutoSet)
    for _, data := range datas {
        sb.Values({{.expressionValues}})
    }
    statement, args := sb.Build()

    var err error
    if session != nil {
        _, err = session.ExecCtx(ctx, statement, args...)
    } else {
        _, err = m.conn.ExecCtx(ctx, statement, args...)
    }
    return err
}

func (m *custom{{.upperStartCamelObject}}Model) FindByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
    sb := sqlbuilder.Select({{.lowerStartCamelObject}}FieldNames...).From(m.table)
	condition.ApplySelect(sb, conds...)
	statement, args := sb.Build()

	var resp []*{{.upperStartCamelObject}}
	var err error

	if session != nil {
		err = session.QueryRowsCtx(ctx, &resp, statement, args...)
	} else {
	    err = m.conn.QueryRowsCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *custom{{.upperStartCamelObject}}Model) FindOneByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) (*{{.upperStartCamelObject}}, error) {
	sb := sqlbuilder.Select({{.lowerStartCamelObject}}FieldNames...).From(m.table)

	condition.ApplySelect(sb, conds...)
	sb.Limit(1)
	statement, args := sb.Build()

	var resp {{.upperStartCamelObject}}
	var err error

	if session != nil {
		err = session.QueryRowCtx(ctx, &resp, statement, args...)
	} else {
	    err = m.conn.QueryRowCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *custom{{.upperStartCamelObject}}Model) PageByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, int64 ,error) {
	sb := sqlbuilder.Select({{.lowerStartCamelObject}}FieldNames...).From(m.table)
	countsb := sqlbuilder.Select("count(*)").From(m.table)

	condition.ApplySelect(sb, conds...)

	var countConds []condition.Condition
    for _, cond := range conds {
    	if cond.Operator != condition.Limit && cond.Operator != condition.Offset {
    		countConds = append(countConds, cond)
    	}
    }
    condition.ApplySelect(countsb, countConds...)

	var resp []*{{.upperStartCamelObject}}
	var err error

	statement, args := sb.Build()

	if session != nil {
		err = session.QueryRowsCtx(ctx, &resp, statement, args...)
	} else {
		err = m.conn.QueryRowsCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, 0, err
	}

    var total int64
    statement, args = countsb.Build()
    if session != nil {
    	err = session.QueryRowCtx(ctx, &total, statement, args...)
    } else {
    	err = m.conn.QueryRowCtx(ctx, &total, statement, args...)
    }
    if err != nil {
    	return nil, 0, err
    }

	return resp, total, nil
}

func (m *custom{{.upperStartCamelObject}}Model) UpdateFieldsByCondition(ctx context.Context, session sqlx.Session, field map[string]any, conds ...condition.Condition) error {
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

	statement, args := sb.Build()

	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, args...)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *custom{{.upperStartCamelObject}}Model) DeleteByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) error {
    if len(conds) == 0 {
		return nil
	}
	sb := sqlbuilder.DeleteFrom(m.table)
	condition.ApplyDelete(sb, conds...)
	statement, args := sb.Build()

	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, args...)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	return err
}