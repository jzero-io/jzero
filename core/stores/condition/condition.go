package condition

import (
	"reflect"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/spf13/cast"
)

type Operator string

func (o Operator) String() string {
	return string(o)
}

type Field string

func (f Field) String() string {
	return string(f)
}

const (
	Equal            Operator = "="
	NotEqual         Operator = "!="
	GreaterThan      Operator = ">"
	LessThan         Operator = "<"
	GreaterEqualThan Operator = ">="
	LessEqualThan    Operator = "<="
	In               Operator = "IN"
	NotIn            Operator = "NOT IN"
	Like             Operator = "LIKE"
	NotLike          Operator = "NOT LIKE"
	Limit            Operator = "LIMIT"
	Offset           Operator = "OFFSET"
	Between          Operator = "BETWEEN"
	NotBetween       Operator = "NOT BETWEEN"
	OrderBy          Operator = "ORDER BY"
	GroupBy          Operator = "GROUP BY"
	Join             Operator = "JOIN"
)

type Condition struct {
	// Skip indicates whether the condition is effective.
	Skip bool

	// SkipFunc The priority is higher than Skip.
	SkipFunc func() bool

	// Or indicates an or condition
	Or bool

	OrOperators  []Operator
	OrFields     []Field
	OrValues     []any
	OrValuesFunc func() []any

	// Field for default and condition
	Field Field

	Operator Operator
	Value    any

	// ValueFunc The priority is higher than Value.
	ValueFunc func() any

	// JoinCondition
	JoinCondition

	WhereClause *sqlbuilder.WhereClause
}

type JoinCondition struct {
	Option sqlbuilder.JoinOption
	Table  string
	OnExpr []string
}

func New(conditions ...Condition) []Condition {
	return conditions
}

func buildExpr(cond *sqlbuilder.Cond, field Field, operator Operator, value any) string {
	switch operator {
	case Equal:
		return cond.Equal(string(field), value)
	case NotEqual:
		return cond.NotEqual(string(field), value)
	case GreaterThan:
		return cond.GreaterThan(string(field), value)
	case LessThan:
		return cond.LessThan(string(field), value)
	case GreaterEqualThan:
		return cond.GreaterEqualThan(string(field), value)
	case LessEqualThan:
		return cond.LessEqualThan(string(field), value)
	case In:
		if len(ToSlice(value)) == 0 {
			// if value is empty, force placeholder nil to avoid sql error
			return cond.In(string(field), nil)
		}
		return cond.In(string(field), ToSlice(value)...)
	case NotIn:
		if len(ToSlice(value)) == 0 {
			// if value is empty, force placeholder nil to avoid sql error
			return cond.NotIn(string(field), nil)
		}
		return cond.NotIn(string(field), ToSlice(value)...)
	case Like:
		return cond.Like(string(field), value)
	case NotLike:
		return cond.NotLike(string(field), value)
	case Between:
		v := ToSlice(value)
		return cond.Between(string(field), v[0], v[1])
	case NotBetween:
		v := ToSlice(value)
		return cond.NotBetween(string(field), v[0], v[1])
	}
	return ""
}

func whereClause(conditions ...Condition) *sqlbuilder.WhereClause {
	clause := sqlbuilder.NewWhereClause()
	cond := sqlbuilder.NewCond()

	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.WhereClause != nil {
			clause.AddWhereClause(c.WhereClause)
			continue
		}
		if c.Or {
			if c.OrValuesFunc != nil {
				c.OrValues = c.OrValuesFunc()
			}
			var expr []string
			for i, field := range c.OrFields {
				if or := buildExpr(cond, field, c.OrOperators[i], c.OrValues[i]); or != "" {
					expr = append(expr, or)
				}
			}
			if len(expr) > 0 {
				clause.AddWhereExpr(cond.Args, cond.Or(expr...))
			}
		} else {
			if c.ValueFunc != nil {
				c.Value = c.ValueFunc()
			}
			if and := buildExpr(cond, c.Field, c.Operator, c.Value); and != "" {
				clause.AddWhereExpr(cond.Args, and)
			}
		}
	}
	return clause
}

func Select(sb sqlbuilder.SelectBuilder, conditions ...Condition) sqlbuilder.SelectBuilder {
	clause := whereClause(conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			sb.Limit(cast.ToInt(c.Value))
		case Offset:
			sb.Offset(cast.ToInt(c.Value))
		case OrderBy:
			sb.OrderBy(cast.ToStringSlice(ToSlice(c.Value))...)
		case GroupBy:
			sb.GroupBy(cast.ToStringSlice(ToSlice(c.Value))...)
		case Join:
			sb.JoinWithOption(c.JoinCondition.Option, c.JoinCondition.Table, cast.ToStringSlice(ToSlice(c.JoinCondition.OnExpr))...)
		}
	}
	if clause != nil {
		sb = *sb.AddWhereClause(clause)
	}
	return sb
}

func Update(builder sqlbuilder.UpdateBuilder, conditions ...Condition) sqlbuilder.UpdateBuilder {
	clause := whereClause(conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			builder.Limit(cast.ToInt(c.Value))
		case OrderBy:
			if len(ToSlice(c.Value)) > 0 {
				builder.OrderBy(cast.ToStringSlice(ToSlice(c.Value))...)
			}
		}
	}
	if clause != nil {
		builder = *builder.AddWhereClause(clause)
	}
	return builder
}

func Delete(builder sqlbuilder.DeleteBuilder, conditions ...Condition) sqlbuilder.DeleteBuilder {
	clause := whereClause(conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			builder.Limit(cast.ToInt(c.Value))
		case OrderBy:
			builder.OrderBy(cast.ToStringSlice(ToSlice(c.Value))...)
		}
	}
	if clause != nil {
		builder = *builder.AddWhereClause(clause)
	}
	return builder
}

func ToSlice(i any) []any {
	if i == nil {
		return []any{}
	}

	switch v := i.(type) {
	case []any:
		return v
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]any, s.Len())
		for i := range a {
			a[i] = s.Index(i).Interface()
		}
		return a
	case
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return []any{i}
	default:
		return []any{}
	}
}
