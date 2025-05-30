var (
{{.lowerStartCamelObject}}FieldNames []string
{{.lowerStartCamelObject}}Rows string
{{.lowerStartCamelObject}}RowsExpectAutoSet string

{{if .withCache}}{{.cacheKeys}}{{end}}
)

func initVars() {
        {{.lowerStartCamelObject}}FieldNames = condition.RawFieldNames(&{{.upperStartCamelObject}}{})
        {{.lowerStartCamelObject}}Rows = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
        {{.lowerStartCamelObject}}RowsExpectAutoSet = strings.Join(condition.RemoveIgnoreColumns({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}", {{end}} {{.ignoreColumns}}), ",")
}
