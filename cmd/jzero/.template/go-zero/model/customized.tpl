func (m *custom{{.upperStartCamelObject}}Model) WithTable(f func(table string) string) {{.lowerStartCamelObject}}Model  {
	mc := &custom{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: m.clone(),
	}
	mc.table = condition.QuoteWithFlavor(m.flavor, f(m.table))
	return mc
}

func (m *custom{{.upperStartCamelObject}}Model) BulkInsert(ctx context.Context, session sqlx.Session, datas []*{{.upperStartCamelObject}}) error {
    if len(datas) == 0 {
        return nil
    }

    sb := sqlbuilder.InsertInto(m.table)
    sb.Cols({{.lowerStartCamelObject}}RowsExpectAutoSet)
    for _, data := range datas {
        sb.Values({{.expressionValues}})
    }
    statement, args := sb.BuildWithFlavor(m.flavor)

    var err error
    if session != nil {
        _, err = session.ExecCtx(ctx, statement, args...)
    } else {
        _, err = m.conn.ExecCtx(ctx, statement, args...)
    }
    return err
}

func (m *custom{{.upperStartCamelObject}}Model) FindSelectedColumnsByCondition(ctx context.Context, session sqlx.Session, columns []string, conditions ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
    if len(columns) == 0 {
        columns = {{.lowerStartCamelObject}}FieldNames
    }

	statement, args := condition.BuildSelectWithFlavor(m.flavor, sqlbuilder.Select(m.withTableColumns(columns...)...).From(m.table), conditions...)

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

func (m *custom{{.upperStartCamelObject}}Model) FindByCondition(ctx context.Context, session sqlx.Session, conditions ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
   return m.FindSelectedColumnsByCondition(ctx, session, {{.lowerStartCamelObject}}FieldNames, conditions...)
}

func (m *custom{{.upperStartCamelObject}}Model) CountByCondition(ctx context.Context, session sqlx.Session, conditions ...condition.Condition) (int64, error) {
   var countconditions []condition.Condition
   for _, cond := range conditions {
    if cond.Operator != condition.Limit && cond.Operator != condition.Offset && cond.Operator != condition.OrderBy && cond.Operator != condition.OrderByDesc && cond.Operator != condition.OrderByAsc {
   	    countconditions = append(countconditions, cond)
    }
   }

   statement, args := condition.BuildSelectWithFlavor(m.flavor, sqlbuilder.Select("count(*)").From(m.table), countconditions...)

   var (
    total int64
    err error
   )

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

func (m *custom{{.upperStartCamelObject}}Model) FindOneByCondition(ctx context.Context, session sqlx.Session, conditions ...condition.Condition) (*{{.upperStartCamelObject}}, error) {
    return m.FindOneSelectedColumnsByCondition(ctx, session, {{.lowerStartCamelObject}}FieldNames, conditions...)
}

func (m *custom{{.upperStartCamelObject}}Model) FindOneSelectedColumnsByCondition(ctx context.Context, session sqlx.Session, columns []string, conditions ...condition.Condition) (*{{.upperStartCamelObject}}, error) {
	statement, args := condition.BuildSelectWithFlavor(m.flavor, sqlbuilder.Select(m.withTableColumns(columns...)...).From(m.table).Limit(1), conditions...)

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

func (m *custom{{.upperStartCamelObject}}Model) PageByCondition(ctx context.Context, session sqlx.Session, conditions ...condition.Condition) ([]*{{.upperStartCamelObject}}, int64 ,error) {
	statement, args := condition.BuildSelectWithFlavor(m.flavor, sqlbuilder.Select(m.withTableColumns({{.lowerStartCamelObject}}FieldNames...)...).From(m.table), conditions...)

	var resp []*{{.upperStartCamelObject}}
	var err error

	if session != nil {
		err = session.QueryRowsCtx(ctx, &resp, statement, args...)
	} else {
		err = m.conn.QueryRowsCtx(ctx, &resp, statement, args...)
	}
	if err != nil {
		return nil, 0, err
	}

	total, err := m.CountByCondition(ctx, session, conditions...)
	if err != nil {
        return nil, 0, err
    }

	return resp, total, nil
}

func (m *custom{{.upperStartCamelObject}}Model) UpdateFieldsByCondition(ctx context.Context, session sqlx.Session, data map[string]any, conditions ...condition.Condition) error {
    if data == nil {
        return nil
    }

	statement, args := condition.BuildUpdateWithFlavor(m.flavor, sqlbuilder.Update(m.table), data, conditions...)

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

func (m *custom{{.upperStartCamelObject}}Model) DeleteByCondition(ctx context.Context, session sqlx.Session, conditions ...condition.Condition) error {
    if len(conditions) == 0 {
		return nil
	}
	statement, args := condition.BuildDeleteWithFlavor(m.flavor, sqlbuilder.DeleteFrom(m.table), conditions...)

	var err error
	if session != nil {
		_, err = session.ExecCtx(ctx, statement, args...)
	} else {
		_, err = m.conn.ExecCtx(ctx, statement, args...)
	}
	return err
}