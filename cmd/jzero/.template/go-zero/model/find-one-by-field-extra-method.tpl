func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", {{.primaryKeyLeft}}, primary)
}

func (m *default{{.upperStartCamelObject}}Model) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
    sb := sqlbuilder.Select({{.lowerStartCamelObject}}Rows).From(m.table)
    sb.Where(sb.EQ(condition.QuoteWithFlavor(m.flavor, "{{.originalPrimaryField}}"), primary))
    sql, args := sb.BuildWithFlavor(m.flavor)
	return conn.QueryRowCtx(ctx, v, sql, args...)
}
