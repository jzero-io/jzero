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

    {{if $.expressionValues}}for _, data := range datas {
            sb.Values({{.expressionValues}})
        }{{end}}

    statement, args := sb.BuildWithFlavor(m.flavor)

    var err error
    if session != nil {
        _, err = session.ExecCtx(ctx, statement, args...)
    } else {
        _, err = m.conn.ExecCtx(ctx, statement, args...)
    }
    return err
}

func (m *custom{{.upperStartCamelObject}}Model) FindFieldsByCondition(ctx context.Context, session sqlx.Session, fields []condition.Field, conditions ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
    if len(fields) == 0 {
        fields = condition.ToFieldSlice({{.lowerStartCamelObject}}FieldNames)
    }

	statement, args := condition.BuildSelectWithFlavor(m.flavor, sqlbuilder.Select(m.withTableFields(cast.ToStringSlice(fields)...)...).From(m.table), conditions...)

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

func (m *custom{{.upperStartCamelObject}}Model) FindSelectedColumnsByCondition(ctx context.Context, session sqlx.Session, columns []string, conditions ...condition.Condition) ([]*{{.upperStartCamelObject}}, error) {
    return m.FindFieldsByCondition(ctx, session, condition.ToFieldSlice(columns), conditions...)
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
    return m.FindOneFieldsByCondition(ctx, session, condition.ToFieldSlice({{.lowerStartCamelObject}}FieldNames), conditions...)
}

func (m *custom{{.upperStartCamelObject}}Model) FindOneFieldsByCondition(ctx context.Context, session sqlx.Session, fields []condition.Field, conditions ...condition.Condition) (*{{.upperStartCamelObject}}, error) {
	statement, args := condition.BuildSelectWithFlavor(m.flavor, sqlbuilder.Select(m.withTableFields(cast.ToStringSlice({{.lowerStartCamelObject}}FieldNames)...)...).From(m.table).Limit(1), conditions...)

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
	statement, args := condition.BuildSelectWithFlavor(m.flavor, sqlbuilder.Select(m.withTableFields({{.lowerStartCamelObject}}FieldNames...)...).From(m.table), conditions...)

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

{{if .withCache}}
func (m *custom{{.upperStartCamelObject}}Model) UpdateFieldsByCondition(ctx context.Context, session sqlx.Session, dataMap map[string]any, conditions ...condition.Condition) error {
    if dataMap == nil {
        return nil
    }

    var fields []condition.Field
    fields = append(fields, {{.data.PrimaryKey.Name.ToCamel}}){{range $ui := .data.UniqueIndex}}
        {{range $uif := $ui}}fields = append(fields, {{$uif.Name.ToCamel}})
        {{end}}
    {{end}}
    fields = lo.Uniq(fields)

    datas, err := m.FindFieldsByCondition(ctx, session, fields, conditions...)
    if err != nil {
    	return err
    }

    var cacheKeys []string

    for _, data := range datas {
        cacheKeys = append(cacheKeys, {{.data.PrimaryCacheKey.DataKeyRight}}){{range $index, $uck := .data.UniqueCacheKey}}
        cacheKeys = append(cacheKeys, {{$uck.DataKeyRight}}){{end}}
    }

	_, err = m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		statement, args := condition.BuildUpdateWithFlavor(m.flavor, sqlbuilder.Update(m.table), dataMap, conditions...)
		if session != nil {
			return session.ExecCtx(ctx, statement, args...)
		}
		return conn.ExecCtx(ctx, statement, args...)
	}, cacheKeys...)
	return err
}
{{else}}
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
{{end}}

{{if .withCache}}
func (m *custom{{.upperStartCamelObject}}Model) DeleteByCondition(ctx context.Context, session sqlx.Session, conditions ...condition.Condition) error {
    if len(conditions) == 0 {
		return nil
	}

    var fields []condition.Field
    fields = append(fields, {{.data.PrimaryKey.Name.ToCamel}}){{range $ui := .data.UniqueIndex}}
        {{range $uif := $ui}}fields = append(fields, {{$uif.Name.ToCamel}})
        {{end}}
    {{end}}
    fields = lo.Uniq(fields)

    datas, err := m.FindFieldsByCondition(ctx, session, fields, conditions...)
    if err != nil {
    	return err
    }

    var cacheKeys []string

    for _, data := range datas {
        cacheKeys = append(cacheKeys, {{.data.PrimaryCacheKey.DataKeyRight}}){{range $index, $uck := .data.UniqueCacheKey}}
        cacheKeys = append(cacheKeys, {{$uck.DataKeyRight}}){{end}}
    }

	_, err = m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		statement, args := condition.BuildDeleteWithFlavor(m.flavor, sqlbuilder.DeleteFrom(m.table), conditions...)
		if session != nil {
			return session.ExecCtx(ctx, statement, args...)
		}
		return conn.ExecCtx(ctx, statement, args...)
	}, cacheKeys...)
	return err
}
{{else}}
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
{{end}}