var (
{{.lowerStartCamelObject}}FieldNames []string
{{.lowerStartCamelObject}}Rows string
{{.lowerStartCamelObject}}RowsExpectAutoSet string

{{.upperStartCamelObject}}Field = struct {
    {{range $index, $v := .data.Fields}}{{$v.Name.ToCamel}} condition.Field
    {{end}}
} {
    {{range $index, $v := .data.Fields}}{{$v.Name.ToCamel}}: "{{$v.NameOriginal}}",
    {{end}}
}

{{if .withCache}}{{.cacheKeys}}{{end}}
)

// Deprecated use {{.upperStartCamelObject}}Field instead
const (
{{range $index, $v := .data.Fields}}{{$v.Name.ToCamel}} condition.Field = "{{$v.NameOriginal}}"
{{end}}
)

func init{{.upperStartCamelObject}}Vars() {
        {{.lowerStartCamelObject}}FieldNames = condition.RawFieldNames(&{{.upperStartCamelObject}}{})
        {{.lowerStartCamelObject}}Rows = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
        {{.lowerStartCamelObject}}RowsExpectAutoSet = strings.Join(condition.RemoveIgnoreColumns({{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}", {{end}} {{.ignoreColumns}}), ",")
}
