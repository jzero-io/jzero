func (m *custom{{.upperStartCamelObject}}Model) BulkInsert(ctx context.Context, session sqlx.Session, datas []*{{.upperStartCamelObject}}) error {
    if len(datas) == 0 {
        return nil
    }

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

func (m *custom{{.upperStartCamelObject}}Model) FindSelectedColumnsByCondition(ctx context.Context, session sqlx.Session, columns []string, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
    if len(columns) == 0 {
        columns = {{.lowerStartCamelObject}}FieldNames
    }
    sb := sqlbuilder.Select(columns...).From(m.table)
	builder := condition.Select(*sb, conds...)
	statement, args := builder.Build()

	var resp []*{{.upperStartCamelObject}}
	var err error

	if session != nil {
		err = session.QueryRowsPartialCtx(ctx, &resp, statement, args...)
	} else {
	    err = m.conn.QueryRowsPartialCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *custom{{.upperStartCamelObject}}Model) FindByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
   return m.FindSelectedColumnsByCondition(ctx, session, {{.lowerStartCamelObject}}FieldNames, conds...)
}

func (m *custom{{.upperStartCamelObject}}Model) CountByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) (int64, error) {
   countsb := sqlbuilder.Select("count(*)").From(m.table)

   var countConds []condition.Condition
   for _, cond := range conds {
    if cond.Operator != condition.Limit && cond.Operator != condition.Offset && cond.Operator != condition.OrderBy {
   	    countConds = append(countConds, cond)
    }
   }
   countBuilder := condition.Select(*countsb, countConds...)

   var (
    total int64
    err error
   )
   statement, args := countBuilder.Build()
   if session != nil {
    err = session.QueryRowCtx(ctx, &total, statement, args...)
   } else {
   	err = m.conn.QueryRowCtx(ctx, &total, statement, args...)
   }
   if err != nil {
    return 0, err
   }
   return total, nil
}

func (m *custom{{.upperStartCamelObject}}Model) FindOneByCondition(ctx context.Context, session sqlx.Session, conds ...condition.Condition) (*{{.upperStartCamelObject}}, error) {
	sb := sqlbuilder.Select({{.lowerStartCamelObject}}FieldNames...).From(m.table)

	builder := condition.Select(*sb, conds...)
	builder.Limit(1)
	statement, args := builder.Build()

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
	builder := condition.Select(*sb, conds...)

	var resp []*{{.upperStartCamelObject}}
	var err error

	statement, args := builder.Build()

	if session != nil {
		err = session.QueryRowsCtx(ctx, &resp, statement, args...)
	} else {
		err = m.conn.QueryRowsCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, 0, err
	}

	total, err := m.CountByCondition(ctx, session, conds...)
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
	builder := condition.Update(*sb, conds...)

	var assigns []string
    for key, value := range field {
        assigns = append(assigns, sb.Assign(key, value))
    }
    builder.Set(assigns...)

	statement, args := builder.Build()

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
	builder := condition.Delete(*sb, conds...)
	statement, args := builder.Build()

	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, args...)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	return err
}