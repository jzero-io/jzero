package condition

import (
	"fmt"
	"reflect"
	"strings"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
)

func SelectByWhereRawSql(sb *sqlbuilder.SelectBuilder, originalField string, args ...any) {
	originalFields := strings.Split(originalField, " and ")
	for i, v := range originalFields {
		field := AdaptField(strings.Split(v, " = ")[0])
		sb.Where(sb.EQ(field, args[i]))
	}
}

const dbTag = "db"

// RawFieldNames converts golang struct field into slice string.
func RawFieldNames(in any) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		tagv := fi.Tag.Get(dbTag)
		switch tagv {
		case "-":
			continue
		case "":
			out = append(out, adapt(fi.Name))
		default:
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if tagv == "-" {
				continue
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}
			out = append(out, adapt(tagv))
		}
	}

	return out
}

func RemoveIgnoreColumns(columns []string, ignoreColumns ...string) []string {
	out := append([]string(nil), columns...)

	for _, ic := range ignoreColumns {
		ic = adapt(ic)
		var n int
		for _, v := range out {
			if v != ic {
				out[n] = v
				n++
			}
		}
		out = out[:n]
	}

	return out
}

func AdaptTable(table string) string {
	return adapt(table)
}

func AdaptField(field string) string {
	return adapt(field)
}

func adapt(str string) string {
	str = Unquote(str)
	return sqlbuilder.DefaultFlavor.Quote(str)
}

func Unquote(s string) string {
	if len(s) == 0 {
		return s
	}
	left := s[0]

	if left == '`' || left == '"' {
		s = s[1:]
	}
	if len(s) == 0 {
		return s
	}
	right := s[len(s)-1]
	if right == '`' || right == '"' {
		s = s[0 : len(s)-1]
	}
	return s
}
