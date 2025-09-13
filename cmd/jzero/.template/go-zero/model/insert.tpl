{{if .withCache}}
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
	statement, args := sqlbuilder.NewInsertBuilder().
                InsertInto(m.table).
                Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
                Values({{.expressionValues}}).Build()
    return m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        if session != nil {
            return session.ExecCtx(ctx, statement, args...)
        }
		return conn.ExecCtx(ctx, statement, args...)
	}, {{.keyValues}}){{end}}
}
{{else}}
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
    statement, args := sqlbuilder.NewInsertBuilder().
            InsertInto(m.table).
            Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
            Values({{.expressionValues}}).Build()
	if session != nil {
       return session.ExecCtx(ctx, statement, args...)
	}
	return m.conn.ExecCtx(ctx, statement, args...)
}
{{end}}

{{if .withCache}}
func (m *default{{.upperStartCamelObject}}Model) InsertV2(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
	statement, args := sqlbuilder.NewInsertBuilder().
                InsertInto(m.table).
                Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
                Values({{.expressionValues}}).Returning("{{.data.Table.PrimaryKey.Name.Source}}").Build()
	var primaryKey {{.data.Table.PrimaryKey.Field.DataType}}
	var err error
	if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL || sqlbuilder.DefaultFlavor == sqlbuilder.SQLite {
        err = m.cachedConn.QueryRowNoCacheCtx(ctx, &primaryKey, statement, args...)
        if err != nil {
            return err
        }
        err = m.cachedConn.DelCacheCtx(ctx, {{.keyValues}})
        if err != nil {
            return err
        }
    } else {
        {{if .data.Table.PrimaryKey.AutoIncrement}}result{{else}}_{{end}}, err := m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
            if session != nil {
                return session.ExecCtx(ctx, statement, args...)
            } else {
                return m.conn.ExecCtx(ctx, statement, args...)
            }
        }, {{.keyValues}})
        if err != nil {
            return err
        }
        {{if .data.Table.PrimaryKey.AutoIncrement}}lastInsertId, err := result.LastInsertId()
        if err != nil {
            return err
        }
        primaryKey = {{.data.Table.PrimaryKey.Field.DataType}}(lastInsertId){{end}}
    }
    if err != nil {
        return err
    }
	data.{{.data.Table.PrimaryKey.Name.ToCamel}} = primaryKey
	return nil{{end}}
}
{{else}}
func (m *default{{.upperStartCamelObject}}Model) InsertV2(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) error {
    statement, args := sqlbuilder.NewInsertBuilder().
            InsertInto(m.table).
            Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
            Values({{.expressionValues}}).Returning("{{.data.Table.PrimaryKey.Name.Source}}").Build()
	var primaryKey {{.data.Table.PrimaryKey.Field.DataType}}
	var err error
	if session != nil {
	    if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL || sqlbuilder.DefaultFlavor == sqlbuilder.SQLite {
	        err = session.QueryRowCtx(ctx, &primaryKey, statement, args...)
	    } else {
	        {{if .data.Table.PrimaryKey.AutoIncrement}}result{{else}}_{{end}}, err := session.ExecCtx(ctx, statement, args...)
	        if err != nil {
	            return err
	        }
	        {{if .data.Table.PrimaryKey.AutoIncrement}}lastInsertId, err := result.LastInsertId()
	        if err != nil {
	            return err
	        }
	        primaryKey = {{.data.Table.PrimaryKey.Field.DataType}}(lastInsertId){{end}}
	    }
	} else {
		if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL || sqlbuilder.DefaultFlavor == sqlbuilder.SQLite {
		    err = m.conn.QueryRowCtx(ctx, &primaryKey, statement, args...)
		} else {
		    {{if .data.Table.PrimaryKey.AutoIncrement}}result{{else}}_{{end}}, err := m.conn.ExecCtx(ctx, statement, args...)
		    if err != nil {
		        return err
		    }
		    {{if .data.Table.PrimaryKey.AutoIncrement}}lastInsertId, err := result.LastInsertId()
		    if err != nil {
		        return err
		    }
		    primaryKey = {{.data.Table.PrimaryKey.Field.DataType}}(lastInsertId){{end}}
		}
	}
	data.{{.data.Table.PrimaryKey.Name.ToCamel}} = primaryKey
    return err
}
{{end}}