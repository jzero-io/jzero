package condition

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
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
	fmt.Println(sql)
	fmt.Println(args)
}

func TestChain2(t *testing.T) {
	sb := sqlbuilder.NewDeleteBuilder()
	sb.DeleteFrom("user")
	chain := NewChain()
	conds := chain.
		Like("name", "%"+"j"+"%").
		In("id", []int{}).
		Build()
	builder := Delete(*sb, conds...)
	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
	fmt.Println(builder.String())
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
	fmt.Println(sql)
	fmt.Println(args)
}
