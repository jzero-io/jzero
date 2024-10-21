// Code generated by jzero. DO NOT EDIT.

package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	{{range $v := .Imports}}"{{$v}}"
	{{end}}
)

type Model struct {
    {{range $v := .TablePackages}}{{$v | FirstUpper | ToCamel}} {{$v}}.{{$v | FirstUpper |ToCamel}}Model
    {{end}}
}

func NewModel(conn sqlx.SqlConn) Model {
	return Model{
         {{range $v := .TablePackages}}{{$v | FirstUpper | ToCamel}}: {{$v}}.New{{ $v | FirstUpper | ToCamel }}Model(conn),
         {{end}}
	}
}