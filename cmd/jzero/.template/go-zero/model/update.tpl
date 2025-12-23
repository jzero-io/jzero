{{if .withCache}}
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, session sqlx.Session, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}) error {
	var err error

	{{if .containsIndexCache}}data, err := m.FindOne(ctx, session, newData.{{.upperStartCamelPrimaryKey}})
	if err != nil {
		return err
	}
    {{end}}{{.keys}}
    _, err = m.cachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
        sb := sqlbuilder.Update(m.table)
        var assigns []string
        {{range $index, $v := .data.Fields}}if slices.Contains({{$.lowerStartCamelObject}}RowsExpectAutoFieldNames, condition.QuoteWithFlavor(m.flavor, "{{$v.Name.Source}}")) {
            assigns = append(assigns, sb.Assign(condition.QuoteWithFlavor(m.flavor, "{{$v.Name.Source}}"), {{if $.containsIndexCache}}newData{{else}}data{{end}}.{{$v.Name.ToCamel}}))
        }
        {{end}}
        sb.Set(assigns...)
        sb.Where(sb.EQ(condition.QuoteWithFlavor(m.flavor, "{{.originalPrimaryKey}}"), {{if $.containsIndexCache}}newData{{else}}data{{end}}.{{.upperStartCamelPrimaryKey}}))
        statement, args := sb.BuildWithFlavor(m.flavor)
        if session != nil {
            return session.ExecCtx(ctx, statement, args...)
        }
		return conn.ExecCtx(ctx, statement, args...)
	}, {{.keyValues}})
	return err
}
{{else}}
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) error {
	sb := sqlbuilder.Update(m.table)
	var assigns []string
    {{range $index, $v := .data.Fields}}if slices.Contains({{$.lowerStartCamelObject}}RowsExpectAutoFieldNames, condition.QuoteWithFlavor(m.flavor, "{{$v.Name.Source}}")) {
        assigns = append(assigns, sb.Assign(condition.QuoteWithFlavor(m.flavor, "{{$v.Name.Source}}"), data.{{$v.Name.ToCamel}}))
    }
    {{end}}
    sb.Set(assigns...)
    sb.Where(sb.EQ(condition.QuoteWithFlavor(m.flavor, "{{.originalPrimaryKey}}"), data.{{.upperStartCamelPrimaryKey}}))
    statement, args := sb.BuildWithFlavor(m.flavor)

    var err error
    if session != nil {
        _, err = session.ExecCtx(ctx, statement, args...)
    } else {
        _, err = m.conn.ExecCtx(ctx, statement, args...)
    }
	return err
}
{{end}}