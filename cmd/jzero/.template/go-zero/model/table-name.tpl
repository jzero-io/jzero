func (m *default{{.upperStartCamelObject}}Model)withTableColumns(columns ...string) []string {
    var withTableColumns []string
    for _, col := range columns {
        if strings.Contains(col, ".") {
            withTableColumns = append(withTableColumns, condition.AdaptField(col))
        } else {
            withTableColumns = append(withTableColumns, m.table + "." + condition.AdaptField(col))
        }
    }
    return withTableColumns
}