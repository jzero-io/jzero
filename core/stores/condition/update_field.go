package condition

import (
	"github.com/eddieowens/opts"
	"github.com/huandu/go-sqlbuilder"
)

type UpdateFieldOperator string

const (
	Incr   UpdateFieldOperator = "INCR"
	Decr   UpdateFieldOperator = "DECR"
	Assign UpdateFieldOperator = "ASSIGN"
	Add    UpdateFieldOperator = "ADD"
	Sub    UpdateFieldOperator = "SUB"
	Mul    UpdateFieldOperator = "MUL"
	Div    UpdateFieldOperator = "DIV"
)

type UpdateFieldChain struct {
	fields []UpdateField
}

type UpdateFieldChainOpts struct {
	Skip      bool
	SkipFunc  func() bool
	ValueFunc func() any
}

func (opts UpdateFieldChainOpts) DefaultOptions() UpdateFieldChainOpts {
	return UpdateFieldChainOpts{}
}

func WithUpdateFieldSkip(skip bool) opts.Opt[UpdateFieldChainOpts] {
	return func(c *UpdateFieldChainOpts) {
		c.Skip = skip
	}
}

func WithUpdateFieldSkipFunc(skipFunc func() bool) opts.Opt[UpdateFieldChainOpts] {
	return func(c *UpdateFieldChainOpts) {
		c.SkipFunc = skipFunc
	}
}

func WithUpdateFieldValueFunc(valueFunc func() any) opts.Opt[UpdateFieldChainOpts] {
	return func(c *UpdateFieldChainOpts) {
		c.ValueFunc = valueFunc
	}
}

type UpdateField struct {
	// Skip indicates whether the UpdateField is effective.
	Skip bool

	// SkipFunc The priority is higher than Skip.
	SkipFunc func() bool

	field    Field
	Operator UpdateFieldOperator

	Value any

	// ValueFunc The priority is higher than Value.
	ValueFunc func() any
}

func (u UpdateFieldChain) addUpdateFieldChain(field Field, operator UpdateFieldOperator, value any, op ...opts.Opt[UpdateFieldChainOpts]) UpdateField {
	o := opts.DefaultApply(op...)

	return UpdateField{
		field:     field,
		Operator:  operator,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	}
}

func NewUpdateFieldChain() UpdateFieldChain {
	return UpdateFieldChain{
		fields: make([]UpdateField, 0),
	}
}

func (u UpdateFieldChain) Assign(field Field, value any, op ...opts.Opt[UpdateFieldChainOpts]) UpdateFieldChain {
	u.fields = append(u.fields, u.addUpdateFieldChain(field, Assign, value, op...))
	return u
}

func (u UpdateFieldChain) Incr(field Field, op ...opts.Opt[UpdateFieldChainOpts]) UpdateFieldChain {
	u.fields = append(u.fields, u.addUpdateFieldChain(field, Incr, nil, op...))
	return u
}

func (u UpdateFieldChain) Decr(field Field, op ...opts.Opt[UpdateFieldChainOpts]) UpdateFieldChain {
	u.fields = append(u.fields, u.addUpdateFieldChain(field, Decr, nil, op...))
	return u
}

func (u UpdateFieldChain) Sub(field Field, value any, op ...opts.Opt[UpdateFieldChainOpts]) UpdateFieldChain {
	u.fields = append(u.fields, u.addUpdateFieldChain(field, Sub, value, op...))
	return u
}

func (u UpdateFieldChain) Mul(field Field, value any, op ...opts.Opt[UpdateFieldChainOpts]) UpdateFieldChain {
	u.fields = append(u.fields, u.addUpdateFieldChain(field, Mul, value, op...))
	return u
}

func (u UpdateFieldChain) Div(field Field, value any, op ...opts.Opt[UpdateFieldChainOpts]) UpdateFieldChain {
	u.fields = append(u.fields, u.addUpdateFieldChain(field, Div, value, op...))
	return u
}

func (u UpdateFieldChain) Add(field Field, value any, op ...opts.Opt[UpdateFieldChainOpts]) UpdateFieldChain {
	u.fields = append(u.fields, u.addUpdateFieldChain(field, Add, value, op...))
	return u
}

func (u UpdateFieldChain) Build() map[string]any {
	fieldMap := make(map[string]any)
	for _, field := range u.fields {
		if field.ValueFunc != nil {
			field.Value = field.ValueFunc()
		}
		if field.SkipFunc != nil {
			field.Skip = field.SkipFunc()
		}
		if field.Skip {
			continue
		}
		fieldMap[string(field.field)] = field
	}
	return fieldMap
}

func SetUpdateFields(builder sqlbuilder.UpdateBuilder, data map[string]any) sqlbuilder.UpdateBuilder {
	return SetUpdateFieldsWithFlavor(sqlbuilder.DefaultFlavor, builder, data)
}

func SetUpdateFieldsWithFlavor(flavor sqlbuilder.Flavor, builder sqlbuilder.UpdateBuilder, data map[string]any) sqlbuilder.UpdateBuilder {
	for key, value := range data {
		if uf, ok := value.(UpdateField); ok {
			switch uf.Operator {
			case Assign:
				builder.SetMore(builder.Assign(QuoteWithFlavor(flavor, key), uf.Value))
			case Incr:
				builder.SetMore(builder.Incr(QuoteWithFlavor(flavor, key)))
			case Decr:
				builder.SetMore(builder.Decr(QuoteWithFlavor(flavor, key)))
			case Div:
				builder.SetMore(builder.Div(QuoteWithFlavor(flavor, key), uf.Value))
			case Add:
				builder.SetMore(builder.Add(QuoteWithFlavor(flavor, key), uf.Value))
			case Mul:
				builder.SetMore(builder.Mul(QuoteWithFlavor(flavor, key), uf.Value))
			case Sub:
				builder.SetMore(builder.Sub(QuoteWithFlavor(flavor, key), uf.Value))
			}
		} else {
			builder.SetMore(builder.Assign(QuoteWithFlavor(flavor, key), value))
		}
	}
	return builder
}
