func (m *default{{.upperStartCamelObject}}Model)withTableColumns(columns ...string) []string {
    var withTableColumns []string
    for _, col := range columns {
        if strings.Contains(col, ".") {
            withTableColumns = append(withTableColumns, condition.QuoteWithFlavor(m.flavor, col))
        } else {
            withTableColumns = append(withTableColumns, m.table + "." + condition.QuoteWithFlavor(m.flavor, col))
        }
    }
    return withTableColumns
}