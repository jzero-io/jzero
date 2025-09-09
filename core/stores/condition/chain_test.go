package condition

import (
	"testing"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	sb := sqlbuilder.NewSelectBuilder().Select("name", "age").From("user")

	chain := NewChain()
	conds := chain.
		Equal("field1", "value1", WithSkip(true)).
		Equal("field2", "value2").
		OrderBy("create_time desc").
		OrderBy("sort desc").
		Build()
	builder := Select(*sb, conds...)

	sql, args := builder.Build()
	assert.Equal(t, `SELECT name, age FROM user WHERE field2 = ? ORDER BY create_time desc, sort desc`, sql)
	assert.Equal(t, []any{"value2"}, args)
}

func TestChainJoin(t *testing.T) {
	sb := sqlbuilder.NewSelectBuilder().Select("user.name", "user.age").From("user")
	chain := NewChain()
	conds := chain.
		Equal("user.field", "value2").
		Join(sqlbuilder.InnerJoin, "user_info", "user.id = user_info.user_id").
		Build()
	builder := Select(*sb, conds...)
	sql, args := builder.Build()
	assert.Equal(t, `SELECT user.name, user.age FROM user INNER JOIN user_info ON user.id = user_info.user_id WHERE user.field = ?`, sql)
	assert.Equal(t, []any{"value2"}, args)
}
