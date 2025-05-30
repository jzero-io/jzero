package condition

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestSelectWithCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

	var values []any
	values = append(values, []int{24, 48}, []int{170, 175})

	cds := New(Condition{
		Field:    "name",
		Operator: Equal,
		Value:    "jaronnie",
	}, Condition{
		Or:          true,
		OrFields:    []string{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
	})

	sb := sqlbuilder.NewSelectBuilder().Select("name", "age", "height").From("user")
	builder := Select(*sb, cds...)

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestUpdateWithCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

	var values []any
	values = append(values, []int{24, 48}, []int{170, 175})

	cds := New(Condition{
		Field:    "name",
		Operator: Equal,
		Value:    "jaronnie",
	}, Condition{
		Or:          true,
		OrFields:    []string{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
	})

	sb := sqlbuilder.NewUpdateBuilder().Update("user")
	builder := Update(*sb, cds...)
	builder.Set(sb.Equal("name", "gocloudcoder"))

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)

	builder = Update(*sb, Condition{
		Field:    "age",
		Operator: Equal,
		Value:    30,
	})
	sql, args = builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestDeleteWithCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

	var values []any
	values = append(values, []int{24, 48}, []int{170, 175})

	cds := New(Condition{
		SkipFunc: func() bool {
			return true
		},
		Field:    "name",
		Operator: Equal,
		Value:    "jaronnie",
		ValueFunc: func() any {
			return "jaronnie2"
		},
	}, Condition{
		Or:          true,
		OrFields:    []string{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
		OrValuesFunc: func() []any {
			return []any{[]int{24, 49}, []int{170, 176}}
		},
	})

	sb := sqlbuilder.NewDeleteBuilder().DeleteFrom("user")
	builder := Delete(*sb, cds...)

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestSqlBuilder(t *testing.T) {
	builder := sqlbuilder.NewSelectBuilder().Select("id", "name").From("user")
	builder.Where(builder.Or(builder.Equal("id", 1), builder.Equal("id", "2")))
	builder.Where(builder.And(builder.Equal("name", "jaronnie")))
	fmt.Println(builder.Build())
}

func TestWhereClause(t *testing.T) {
	var values []any
	values = append(values, []int{24, 48}, []int{170, 175})
	cds := New(Condition{
		SkipFunc: func() bool {
			return true
		},
		Field:    "name",
		Operator: Equal,
		Value:    "jaronnie",
		ValueFunc: func() any {
			return "jaronnie2"
		},
	}, Condition{
		Or:          true,
		OrFields:    []string{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
		OrValuesFunc: func() []any {
			return []any{[]int{24, 49}, []int{170, 176}}
		},
	})
	clause := whereClause(cds...)
	statement, args := clause.Build()
	fmt.Println(statement)
	fmt.Println(args)
}

func TestRawWhereClause(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

	rawWhereClause := sqlbuilder.NewWhereClause()

	cond := sqlbuilder.NewCond()
	rawWhereClause.AddWhereExpr(cond.Args,
		cond.Or(
			cond.And(
				cond.EQ("a", 1),
				cond.EQ("b", 2),
			),
			cond.And(
				cond.EQ("c", 3),
				cond.EQ("d", 4),
			),
		))

	cds := New(Condition{
		WhereClause: rawWhereClause,
	}, Condition{
		Field:    "field_with_jzero",
		Value:    123,
		Operator: Equal,
	})

	sb := sqlbuilder.NewSelectBuilder().Select("name", "age", "height").From("user")
	builder := Select(*sb, cds...)

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestChain_In(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		sqlbuilder.DefaultFlavor = sqlbuilder.MySQL
		sb := sqlbuilder.NewSelectBuilder().Select("id").From("users")
		builder := Select(*sb, Condition{
			Field:    "id",
			Operator: In,
			Value:    []int{1},
		})

		sql, args := builder.Build()
		fmt.Println(sql, args)
	})

	t.Run("test2", func(t *testing.T) {
		sqlbuilder.DefaultFlavor = sqlbuilder.MySQL
		sb := sqlbuilder.NewSelectBuilder().Select("id").From("users")
		builder := Select(*sb, Condition{
			Field:    "id",
			Operator: In,
			Value:    []int{},
		})

		sql, args := builder.Build()
		fmt.Println(sql, args)
	})
}
