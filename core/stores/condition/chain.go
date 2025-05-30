package condition

import (
	"github.com/eddieowens/opts"
	"github.com/huandu/go-sqlbuilder"
)

type Chain struct {
	conditions []Condition
}

type ChainOperatorOpts struct {
	Skip     bool
	SkipFunc func() bool

	ValueFunc func() any

	OrValuesFunc func() []any
}

func (opts ChainOperatorOpts) DefaultOptions() ChainOperatorOpts {
	return ChainOperatorOpts{}
}

func WithSkip(skip bool) opts.Opt[ChainOperatorOpts] {
	return func(c *ChainOperatorOpts) {
		c.Skip = skip
	}
}

func WithSkipFunc(skipFunc func() bool) opts.Opt[ChainOperatorOpts] {
	return func(c *ChainOperatorOpts) {
		c.SkipFunc = skipFunc
	}
}

func WithValueFunc(valueFunc func() any) opts.Opt[ChainOperatorOpts] {
	return func(c *ChainOperatorOpts) {
		c.ValueFunc = valueFunc
	}
}

func WithOrValuesFunc(valueFunc func() []any) opts.Opt[ChainOperatorOpts] {
	return func(c *ChainOperatorOpts) {
		c.OrValuesFunc = valueFunc
	}
}

func NewChain(conditions ...Condition) Chain {
	return Chain{conditions: conditions}
}

// Deprecated: Use NewChain instead
func NewChainWithConditions(conditions ...Condition) Chain {
	return Chain{conditions: conditions}
}

func (c Chain) AddCondition(condition Condition) Chain {
	c.conditions = append(c.conditions, condition)
	return c
}

func (c Chain) addChain(field string, operator Operator, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  operator,
		Value:     value,
		ValueFunc: o.ValueFunc,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
	})
	return c
}

func (c Chain) Equal(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, Equal, value, op...)
}

func (c Chain) NotEqual(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, NotEqual, value, op...)
}

func (c Chain) GreaterThan(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, GreaterThan, value, op...)
}

func (c Chain) LessThan(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, LessThan, value, op...)
}

func (c Chain) GreaterEqualThan(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, GreaterEqualThan, value, op...)
}

func (c Chain) LessEqualThan(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, LessEqualThan, value, op...)
}

func (c Chain) Like(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, Like, value, op...)
}

func (c Chain) NotLike(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, NotLike, value, op...)
}

func (c Chain) In(field string, values any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, In, values, op...)
}

func (c Chain) NotIn(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, NotIn, value, op...)
}

func (c Chain) Between(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain(field, Between, value, op...)
}

func (c Chain) Or(fields []string, operators []Operator, values []any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Or:           true,
		OrFields:     fields,
		OrOperators:  operators,
		OrValues:     values,
		Skip:         o.Skip,
		SkipFunc:     o.SkipFunc,
		OrValuesFunc: o.OrValuesFunc,
	})
	return c
}

func (c Chain) OrderBy(value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain("", OrderBy, value, op...)
}

func (c Chain) Limit(value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain("", Limit, value, op...)
}

func (c Chain) Offset(value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain("", Offset, value, op...)
}

func (c Chain) Page(page, pageSize int, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain("", Offset, (page-1)*pageSize, op...).addChain("", Limit, pageSize, op...)
}

func (c Chain) GroupBy(value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	return c.addChain("", GroupBy, value, op...)
}

func (c Chain) Join(option sqlbuilder.JoinOption, table string, onExpr ...string) Chain {
	c.conditions = append(c.conditions, Condition{
		Operator: Join,
		JoinCondition: JoinCondition{
			Table:  table,
			OnExpr: onExpr,
			Option: option,
		},
	})
	return c
}

func (c Chain) Build() []Condition {
	return c.conditions
}

func (c Chain) WhereClause(whereClause *sqlbuilder.WhereClause) Chain {
	c.conditions = append(c.conditions, Condition{
		WhereClause: whereClause,
	})
	return c
}
