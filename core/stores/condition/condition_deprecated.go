package condition

import (
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/spf13/cast"
)

// AdaptTable: return quote table
// Deprecated use Quote instead
func AdaptTable(table string) string {
	return QuoteWithFlavor(sqlbuilder.DefaultFlavor, table)
}

// AdaptField: return quote field
// Deprecated use Quote instead
func AdaptField(field string) string {
	return QuoteWithFlavor(sqlbuilder.DefaultFlavor, field)
}

// Select return select with condition builder
// Deprecated: Use BuildSelect instead
func Select(builder sqlbuilder.SelectBuilder, conditions ...Condition) sqlbuilder.SelectBuilder {
	return SelectWithFlavor(sqlbuilder.DefaultFlavor, builder, conditions...)
}

// SelectWithFlavor return flavor select with condition builder
// Deprecated: Use BuildSelectWithFlavor instead
func SelectWithFlavor(flavor sqlbuilder.Flavor, builder sqlbuilder.SelectBuilder, conditions ...Condition) sqlbuilder.SelectBuilder {
	builder.SetFlavor(flavor)
	clause := whereClause(flavor, conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			builder.Limit(cast.ToInt(c.Value))
		case Offset:
			builder.Offset(cast.ToInt(c.Value))
		case OrderBy:
			builder.OrderBy(cast.ToStringSlice(ToSlice(c.Value))...)
		case OrderByDesc:
			builder.OrderByDesc(QuoteWithFlavor(flavor, string(c.Field)))
		case OrderByAsc:
			builder.OrderByAsc(QuoteWithFlavor(flavor, string(c.Field)))
		case GroupBy:
			// compatibility with old version
			if c.Value != nil {
				builder.GroupBy(cast.ToStringSlice(ToSlice(c.Value))...)
			} else {
				builder.GroupBy(QuoteWithFlavor(flavor, string(c.Field)))
			}
		case Join:
			builder.JoinWithOption(c.JoinCondition.Option, c.JoinCondition.Table, cast.ToStringSlice(ToSlice(c.JoinCondition.OnExpr))...)
		}
	}
	if clause != nil {
		builder = *builder.AddWhereClause(clause)
	}
	return builder
}

// Update return update with condition builder
// Deprecated: Use BuildUpdate instead
func Update(builder sqlbuilder.UpdateBuilder, conditions ...Condition) sqlbuilder.UpdateBuilder {
	return UpdateWithFlavor(sqlbuilder.DefaultFlavor, builder, conditions...)
}

// UpdateWithFlavor return flavor update with condition builder
// Deprecated: Use BuildUpdateWithFlavor instead
func UpdateWithFlavor(flavor sqlbuilder.Flavor, builder sqlbuilder.UpdateBuilder, conditions ...Condition) sqlbuilder.UpdateBuilder {
	builder.SetFlavor(flavor)
	clause := whereClause(flavor, conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			builder.Limit(cast.ToInt(c.Value))
		case OrderBy:
			builder.OrderBy(cast.ToStringSlice(ToSlice(c.Value))...)
		case OrderByDesc:
			builder.OrderByDesc(QuoteWithFlavor(flavor, string(c.Field)))
		case OrderByAsc:
			builder.OrderByAsc(QuoteWithFlavor(flavor, string(c.Field)))
		}
	}
	if clause != nil {
		builder = *builder.AddWhereClause(clause)
	}
	return builder
}

// Delete return delete with condition builder
// Deprecated: Use BuildUDelete instead
func Delete(builder sqlbuilder.DeleteBuilder, conditions ...Condition) sqlbuilder.DeleteBuilder {
	return DeleteWithFlavor(sqlbuilder.DefaultFlavor, builder, conditions...)
}

// DeleteWithFlavor return delete with condition builder
// Deprecated: Use BuildUDeleteWithFlavor instead
func DeleteWithFlavor(flavor sqlbuilder.Flavor, builder sqlbuilder.DeleteBuilder, conditions ...Condition) sqlbuilder.DeleteBuilder {
	builder.SetFlavor(flavor)
	clause := whereClause(flavor, conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			builder.Limit(cast.ToInt(c.Value))
		case OrderBy:
			builder.OrderBy(cast.ToStringSlice(ToSlice(c.Value))...)
		case OrderByDesc:
			builder.OrderByDesc(QuoteWithFlavor(flavor, string(c.Field)))
		case OrderByAsc:
			builder.OrderByAsc(QuoteWithFlavor(flavor, string(c.Field)))
		}
	}
	if clause != nil {
		builder = *builder.AddWhereClause(clause)
	}
	return builder
}
