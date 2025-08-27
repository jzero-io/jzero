package model

import (
	"github.com/eddieowens/opts"
	"github.com/jzero-io/jzero/core/stores/monx"
)

const {{.Type}}CollectionName = "{{.snakeType}}"

var _ {{.Type}}Model = (*custom{{.Type}}Model)(nil)

type (
    // {{.Type}}Model is an interface to be customized, add more methods here,
    // and implement the added methods in custom{{.Type}}Model.
    {{.Type}}Model interface {
        {{.lowerType}}Model
    }

    custom{{.Type}}Model struct {
        *default{{.Type}}Model
    }
)

// New{{.Type}}Model returns a model for the mongo.
func New{{.Type}}Model(url, db string, op ...opts.Opt[monx.MonOpts]) {{.Type}}Model {
    return &custom{{.Type}}Model{
    	default{{.Type}}Model: newDefault{{.Type}}Model(url, db, op...),
    }
}