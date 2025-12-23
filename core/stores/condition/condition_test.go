package condition

import (
	"testing"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
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
		OrFields:    []Field{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
	})

	sb := sqlbuilder.NewSelectBuilder().Select("name", "age", "height").From("user")
	sql, args := BuildSelect(sb, cds...)

	assert.Equal(t, `SELECT name, age, height FROM user WHERE `+"`name`"+` = ? AND (`+"`age`"+` BETWEEN ? AND ? OR `+"`height`"+` BETWEEN ? AND ?)`, sql)
	assert.Equal(t, []any{"jaronnie", 24, 48, 170, 175}, args)
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
		OrFields:    []Field{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
	})

	sb := sqlbuilder.NewUpdateBuilder().Update("user")
	sql, args := BuildUpdate(sb.Clone(), map[string]any{}, cds...)

	assert.Equal(t, `UPDATE user WHERE `+"`name`"+` = ? AND (`+"`age`"+` BETWEEN ? AND ? OR `+"`height`"+` BETWEEN ? AND ?)`, sql)
	assert.Equal(t, []any{"jaronnie", 24, 48, 170, 175}, args)

	sql, args = BuildUpdate(sb.Clone(), map[string]any{}, Condition{
		Field:    "age",
		Operator: Equal,
		Value:    30,
	})
	assert.Equal(t, `UPDATE user WHERE `+"`age`"+` = ?`, sql)
	assert.Equal(t, []any{30}, args)
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
		OrFields:    []Field{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
		OrValuesFunc: func() []any {
			return []any{[]int{24, 49}, []int{170, 176}}
		},
	})

	sb := sqlbuilder.NewDeleteBuilder().DeleteFrom("user")
	sql, args := BuildDelete(sb, cds...)

	assert.Equal(t, `DELETE FROM user WHERE (`+"`age`"+` BETWEEN ? AND ? OR `+"`height`"+` BETWEEN ? AND ?)`, sql)
	assert.Equal(t, []any{24, 49, 170, 176}, args)
}

func TestSqlBuilder(t *testing.T) {
	builder := sqlbuilder.NewSelectBuilder().Select("id", "name").From("user")
	builder.Where(builder.Or(builder.Equal("id", 1), builder.Equal("id", 2)))
	builder.Where(builder.And(builder.Equal("name", "jaronnie")))
	sql, args := builder.Build()
	assert.Equal(t, "SELECT id, name FROM user WHERE (id = ? OR id = ?) AND (name = ?)", sql)
	assert.Equal(t, []any{1, 2, "jaronnie"}, args)
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
		OrFields:    []Field{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
		OrValuesFunc: func() []any {
			return []any{[]int{24, 49}, []int{170, 176}}
		},
	})
	clause := whereClause(sqlbuilder.DefaultFlavor, cds...)
	statement, args := clause.Build()
	assert.Equal(t, "WHERE (`age` BETWEEN ? AND ? OR `height` BETWEEN ? AND ?)", statement)
	assert.Equal(t, []any{24, 49, 170, 176}, args)
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
	sql, args := BuildSelect(sb, cds...)
	assert.Equal(t, "SELECT name, age, height FROM user WHERE ((a = ? AND b = ?) OR (c = ? AND d = ?)) AND `field_with_jzero` = ?", sql)
	assert.Equal(t, []any{1, 2, 3, 4, 123}, args)
}

func TestChain_In(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		sqlbuilder.DefaultFlavor = sqlbuilder.MySQL
		sb := sqlbuilder.NewSelectBuilder().Select("id").From("users")
		sql, args := BuildSelect(sb, Condition{
			Field:    "id",
			Operator: In,
			Value:    []int{1},
		})

		assert.Equal(t, "SELECT id FROM users WHERE `id` IN (?)", sql)
		assert.Equal(t, []any{1}, args)
	})

	t.Run("test2", func(t *testing.T) {
		sqlbuilder.DefaultFlavor = sqlbuilder.MySQL
		sb := sqlbuilder.NewSelectBuilder().Select("id").From("users")
		sql, args := BuildSelect(sb, Condition{
			Field:    "id",
			Operator: In,
			Value:    []int{},
		})

		assert.Equal(t, "SELECT id FROM users WHERE `id` IN (?)", sql)
		assert.Equal(t, []any{nil}, args)
	})
}

func TestToFieldSlice(t *testing.T) {
	slice := ToFieldSlice([]string{"age", "height"})
	assert.Equal(t, []Field{"age", "height"}, slice)
}
