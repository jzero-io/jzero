package condition

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestSetUpdateFields(t *testing.T) {
	builder := SetUpdateFieldsWithFlavor(sqlbuilder.MySQL, *sqlbuilder.NewUpdateBuilder().Update("users"), map[string]any{
		"age": UpdateField{
			Operator: Add,
			Value:    12,
		},
		"version": UpdateField{
			Operator: Incr,
		},
		"name": "jaronnie",
	})

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestSetUpdateFieldChain(t *testing.T) {
	builder := SetUpdateFieldsWithFlavor(sqlbuilder.MySQL, *sqlbuilder.NewUpdateBuilder().Update("users"), NewUpdateFieldChain().
		Assign("name", "jaronnie", WithUpdateFieldSkip(true)).
		Incr("version").
		Add("age", 12, WithUpdateFieldValueFunc(func() any {
			return 15
		})).
		Build())

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
}
