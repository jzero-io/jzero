package condition

import (
	"testing"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	sb := sqlbuilder.NewSelectBuilder().Select("name", "age").From("user")

	chain := NewChain().
		Equal("field1", "value1", WithSkip(true)).
		Equal("field2", "value2").
		OrderByDesc("create_time").
		OrderByAsc("sort")
	statement, args := BuildSelectWithFlavor(sqlbuilder.MySQL, sb, chain.Build()...)

	assert.Equal(t, "SELECT name, age FROM user WHERE `field2` = ? ORDER BY `create_time` DESC, `sort` ASC", statement)
	assert.Equal(t, []any{"value2"}, args)
}

func TestChainJoin(t *testing.T) {
	sb := sqlbuilder.NewSelectBuilder().Select("user.name", "user.age").From("user")
	chain := NewChain().
		Equal("user.field", "value2").
		Join(sqlbuilder.InnerJoin, "user_info", "user.id = user_info.user_id")

	statement, args := BuildSelectWithFlavor(sqlbuilder.MySQL, sb, chain.Build()...)
	assert.Equal(t, "SELECT user.name, user.age FROM user INNER JOIN user_info ON user.id = user_info.user_id WHERE `user`.`field` = ?", statement)
	assert.Equal(t, []any{"value2"}, args)
}

func TestChainIsNull(t *testing.T) {
	sb := sqlbuilder.NewSelectBuilder().Select("user.name", "user.age").From("user")
	chain := NewChain().
		Equal("user.field", "value2").
		IsNull("delete_at")

	statement, args := BuildSelectWithFlavor(sqlbuilder.MySQL, sb, chain.Build()...)
	assert.Equal(t, "SELECT user.name, user.age FROM user WHERE `user`.`field` = ? AND `delete_at` IS NULL", statement)
	assert.Equal(t, []any{"value2"}, args)
}

func TestChainIsNotNull(t *testing.T) {
	sb := sqlbuilder.NewSelectBuilder().Select("user.name", "user.age").From("user")
	chain := NewChain().
		Equal("user.field", "value2").
		IsNotNull("delete_at")

	statement, args := BuildSelectWithFlavor(sqlbuilder.MySQL, sb, chain.Build()...)
	assert.Equal(t, "SELECT user.name, user.age FROM user WHERE `user`.`field` = ? AND `delete_at` IS NOT NULL", statement)
	assert.Equal(t, []any{"value2"}, args)
}

func TestSelectForUpdate(t *testing.T) {
	sb := sqlbuilder.NewSelectBuilder().Select("name", "age").From("user")
	chain := NewChain().
		Equal("field1", "value1", WithSkip(true)).
		Equal("field2", "value2").
		OrderByDesc("create_time").
		OrderByAsc("sort").
		ForUpdate()

	statement, args := BuildSelectWithFlavor(sqlbuilder.MySQL, sb, chain.Build()...)
	assert.Equal(t, "SELECT name, age FROM user WHERE `field2` = ? ORDER BY `create_time` DESC, `sort` ASC FOR UPDATE", statement)
	assert.Equal(t, []any{"value2"}, args)
}
