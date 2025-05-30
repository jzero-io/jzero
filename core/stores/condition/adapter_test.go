package condition

import (
	"testing"

	"github.com/huandu/go-sqlbuilder"
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
			t.Log(sql, arguments)
		})
	}
}
