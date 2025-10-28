package condition

import (
	"testing"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestSelectByWhereRawSql(t *testing.T) {
	type args struct {
		sb            *sqlbuilder.SelectBuilder
		originalField string
		args          []any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				sb:            sqlbuilder.NewSelectBuilder(),
				originalField: "`sys_user_id` = ? and `sys_authority_authority_id` = ?",
				args:          []any{1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
			SelectByWhereRawSql(tt.args.sb, tt.args.originalField, tt.args.args...)

			sql, arguments := tt.args.sb.Build()
			assert.Equal(t, `WHERE "sys_user_id" = ? AND "sys_authority_authority_id" = ?`, sql)
			assert.Equal(t, []any{1, 1}, arguments)
		})
	}
}

func TestQuoteWithFlavor(t *testing.T) {
	type args struct {
		flavor sqlbuilder.Flavor
		str    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				flavor: sqlbuilder.MySQL,
				str:    "manage_role.id",
			},
			want: "`manage_role`.`id`",
		},
		{
			name: "test2",
			args: args{
				flavor: sqlbuilder.MySQL,
				str:    "id",
			},
			want: "`id`",
		},
		{
			name: "test3",
			args: args{
				flavor: sqlbuilder.MySQL,
				str:    "`manage_role`.`id`",
			},
			want: "`manage_role`.`id`",
		},
		{
			name: "test4",
			args: args{
				flavor: sqlbuilder.PostgreSQL,
				str:    "`manage_role`.`id`",
			},
			want: `"manage_role"."id"`,
		},
		{
			name: "test4",
			args: args{
				flavor: sqlbuilder.PostgreSQL,
				str:    "id",
			},
			want: `"id"`,
		},
		{
			name: "test5",
			args: args{
				flavor: sqlbuilder.PostgreSQL,
				str:    `"manage_role"."id"`,
			},
			want: `"manage_role"."id"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, QuoteWithFlavor(tt.args.flavor, tt.args.str), "QuoteWithFlavor(%v, %v)", tt.args.flavor, tt.args.str)
		})
	}
}
