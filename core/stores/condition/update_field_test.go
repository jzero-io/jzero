package condition

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestUpdateFields(t *testing.T) {
	sql, args := BuildUpdateWithFlavor(sqlbuilder.MySQL, sqlbuilder.NewUpdateBuilder().Update("users"), map[string]any{
		"age": UpdateField{
			Operator: Add,
			Value:    12,
		},
		"version": UpdateField{
			Operator: Incr,
		},
		"name": "jaronnie",
	}, NewChain().Equal("id", 1).Build()...)

	fmt.Println(sql)
	fmt.Println(args)
}

func TestUpdateFieldChain(t *testing.T) {
	sql, args := BuildUpdateWithFlavor(sqlbuilder.MySQL, sqlbuilder.NewUpdateBuilder().Update("users"), NewUpdateFieldChain().
		Assign("name", "jaronnie", WithUpdateFieldSkip(true)).
		Incr("version").
		Add("age", 12, WithUpdateFieldValueFunc(func() any {
			return 15
		})).
		Build(), NewChain().Equal("id", 1).Build()...)

	fmt.Println(sql)
	fmt.Println(args)
}
