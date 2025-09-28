func (m *default{{.upperStartCamelObject}}Model) TableName() string {
	return m.table
}

func (m *default{{.upperStartCamelObject}}Model)withTableColumns(columns ...string) []string {
    var withTableColumns []string
    for _, col := range columns {
        if strings.Contains(col, ".") {
            withTableColumns = append(withTableColumns, col)
        } else {
            withTableColumns = append(withTableColumns, m.table + "." + col)
        }
    }
    return withTableColumns
}