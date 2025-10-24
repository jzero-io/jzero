{{if .withCache}}
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
	statement, args := sqlbuilder.NewInsertBuilder().
                InsertInto(m.table).
                Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
                Values({{.expressionValues}}).BuildWithFlavor(m.flavor)
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
            Values({{.expressionValues}}).BuildWithFlavor(m.flavor)
	if session != nil {
       return session.ExecCtx(ctx, statement, args...)
	}
	return m.conn.ExecCtx(ctx, statement, args...)
}
{{end}}

{{if .withCache}}
func (m *default{{.upperStartCamelObject}}Model) InsertV2(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
	var statement string
	var args []any
	{{if .data.Table.PrimaryKey.AutoIncrement}}if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL || sqlbuilder.DefaultFlavor == sqlbuilder.SQLite {
		statement, args = sqlbuilder.NewInsertBuilder().
			InsertInto(m.table).
			Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
			Values({{.expressionValues}}).Returning("{{.data.Table.PrimaryKey.Name.Source}}").BuildWithFlavor(m.flavor)
	} else {
		statement, args = sqlbuilder.NewInsertBuilder().
			InsertInto(m.table).
			Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
			Values({{.expressionValues}}).BuildWithFlavor(m.flavor)
	}{{else}}statement, args = sqlbuilder.NewInsertBuilder().
		InsertInto(m.table).
		Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
		Values({{.expressionValues}}).BuildWithFlavor(m.flavor){{end}}
	{{if .data.Table.PrimaryKey.AutoIncrement}}var primaryKey {{.data.Table.PrimaryKey.Field.DataType}}{{end}}
	var err error
	{{if .data.Table.PrimaryKey.AutoIncrement}}if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL || sqlbuilder.DefaultFlavor == sqlbuilder.SQLite {
        if session != nil {
        	err = session.QueryRowCtx(ctx, &primaryKey, statement, args...)
        	if err != nil {
        		return err
        	}
        } else {
        	err = m.cachedConn.QueryRowNoCacheCtx(ctx, &primaryKey, statement, args...)
        	if err != nil {
        		return err
        	}
        }
        err = m.cachedConn.DelCacheCtx(ctx, {{.keyValues}})
        if err != nil {
        	return err
        }
    } else {
        result, err := m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
            if session != nil {
                return session.ExecCtx(ctx, statement, args...)
            } else {
                return conn.ExecCtx(ctx, statement, args...)
            }
        }, {{.keyValues}})
        if err != nil {
            return err
        }
        lastInsertId, err := result.LastInsertId()
        if err != nil {
            return err
        }
        primaryKey = {{.data.Table.PrimaryKey.Field.DataType}}(lastInsertId)
    }
	data.{{.data.Table.PrimaryKey.Name.ToCamel}} = primaryKey{{else}}_, err = m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        if session != nil {
            return session.ExecCtx(ctx, statement, args...)
        } else {
            return conn.ExecCtx(ctx, statement, args...)
        }
    }, {{.keyValues}})
    if err != nil {
        return err
    }{{end}}
	return nil{{end}}
}
{{else}}
func (m *default{{.upperStartCamelObject}}Model) InsertV2(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) error {
    var statement string
    var args []any
    {{if .data.Table.PrimaryKey.AutoIncrement}}if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL || sqlbuilder.DefaultFlavor == sqlbuilder.SQLite {
        statement, args = sqlbuilder.NewInsertBuilder().
                InsertInto(m.table).
                Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
                Values({{.expressionValues}}).Returning("{{.data.Table.PrimaryKey.Name.Source}}").BuildWithFlavor(m.flavor)
    } else {
        statement, args = sqlbuilder.NewInsertBuilder().
                InsertInto(m.table).
                Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
                Values({{.expressionValues}}).BuildWithFlavor(m.flavor)
    }{{else}}statement, args = sqlbuilder.NewInsertBuilder().
            InsertInto(m.table).
            Cols({{.lowerStartCamelObject}}RowsExpectAutoSet).
            Values({{.expressionValues}}).BuildWithFlavor(m.flavor){{end}}
	{{if .data.Table.PrimaryKey.AutoIncrement}}var primaryKey {{.data.Table.PrimaryKey.Field.DataType}}{{end}}
	var err error
	{{if .data.Table.PrimaryKey.AutoIncrement}}if session != nil {
	    if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL || sqlbuilder.DefaultFlavor == sqlbuilder.SQLite {
	        err = session.QueryRowCtx(ctx, &primaryKey, statement, args...)
	        if err != nil {
	            return err
	        }
	    } else {
	        result, err := session.ExecCtx(ctx, statement, args...)
	        if err != nil {
	            return err
	        }
	        lastInsertId, err := result.LastInsertId()
	        if err != nil {
	            return err
	        }
	        primaryKey = {{.data.Table.PrimaryKey.Field.DataType}}(lastInsertId)
	    }
	} else {
		if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL || sqlbuilder.DefaultFlavor == sqlbuilder.SQLite {
		    err = m.conn.QueryRowCtx(ctx, &primaryKey, statement, args...)
		} else {
		    result, err := m.conn.ExecCtx(ctx, statement, args...)
		    if err != nil {
		        return err
		    }
		    lastInsertId, err := result.LastInsertId()
		    if err != nil {
		        return err
		    }
		    primaryKey = {{.data.Table.PrimaryKey.Field.DataType}}(lastInsertId)
		}
	}
	data.{{.data.Table.PrimaryKey.Name.ToCamel}} = primaryKey{{else}}if session != nil {
	    _, err = session.ExecCtx(ctx, statement, args...)
	} else {
	    _, err = m.conn.ExecCtx(ctx, statement, args...)
	}{{end}}
    return err
}
{{end}}