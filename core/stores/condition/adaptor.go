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
		field := strings.Split(v, " = ")[0]
		if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL {
			field = Unquote(field)
			field = fmt.Sprintf(`"%s"`, field)
		}
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
		// gets us a StructField
		fi := typ.Field(i)
		tagv := fi.Tag.Get(dbTag)
		switch tagv {
		case "-":
			continue
		case "":
			if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL {
				out = append(out, fmt.Sprintf(`"%s"`, fi.Name))
			} else {
				out = append(out, fmt.Sprintf("`%s`", fi.Name))
			}
		default:
			// get tag name with the tag option, e.g.:
			// `db:"id"`
			// `db:"id,type=char,length=16"`
			// `db:",type=char,length=16"`
			// `db:"-,type=char,length=16"`
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if tagv == "-" {
				continue
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}
			if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL {
				out = append(out, fmt.Sprintf(`"%s"`, tagv))
			} else {
				out = append(out, fmt.Sprintf("`%s`", tagv))
			}
		}
	}

	return out
}

func RemoveIgnoreColumns(strings []string, strs ...string) []string {
	out := append([]string(nil), strings...)

	for _, str := range strs {
		if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL {
			str = fmt.Sprintf(`"%s"`, Unquote(str))
		}

		var n int
		for _, v := range out {
			if v != str {
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
	switch sqlbuilder.DefaultFlavor {
	case sqlbuilder.PostgreSQL:
		str = Unquote(str)
		return fmt.Sprintf(`"%s"`, str)
	default:
		return str
	}
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
