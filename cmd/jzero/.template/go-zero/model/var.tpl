var (
    {{.lowerStartCamelObject}}FieldNames []string
    {{.lowerStartCamelObject}}Rows string
    {{.lowerStartCamelObject}}RowsExpectAutoFieldNames []string
    {{.lowerStartCamelObject}}RowsExpectAutoSet string

    {{if .withCache}}{{.cacheKeys}}{{end}}
)

const (
    {{range $index, $v := .data.Fields}}{{$v.Name.ToCamel}} condition.Field = "{{$v.NameOriginal}}"
    {{end}}
)

func init{{.upperStartCamelObject}}Vars(flavor sqlbuilder.Flavor) {
    {{.lowerStartCamelObject}}FieldNames = condition.RawFieldNamesWithFlavor(flavor, &{{.upperStartCamelObject}}{})
    {{.lowerStartCamelObject}}Rows = strings.Join({{.lowerStartCamelObject}}FieldNames, ",")
    {{.lowerStartCamelObject}}RowsExpectAutoFieldNames = condition.RemoveIgnoreColumnsWithFlavor(flavor, {{.lowerStartCamelObject}}FieldNames, {{if .autoIncrement}}"{{.originalPrimaryKey}}", {{end}} {{.ignoreColumns}})
    {{.lowerStartCamelObject}}RowsExpectAutoSet = strings.Join({{.lowerStartCamelObject}}RowsExpectAutoFieldNames, ",")
}
