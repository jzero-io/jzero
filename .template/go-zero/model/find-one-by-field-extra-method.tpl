func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", {{.primaryKeyLeft}}, primary)
}

func (m *default{{.upperStartCamelObject}}Model) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
    sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
    sb.Where(sb.EQ("{{.originalPrimaryField}}", primary))
    sql, args := sb.Build()
	return conn.QueryRowCtx(ctx, v, sql, args...)
}
