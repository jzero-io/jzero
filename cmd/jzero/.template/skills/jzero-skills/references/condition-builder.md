# Condition Builder

## Overview

The `condition` package provides a fluent, type-safe way to build database query conditions. It's used with all `*ByCondition` generated methods and supports complex queries with a clean API.

## Import

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"  // Import model for field constants
)
```

## Use Generated Field Constants

**Always use generated field constants instead of hardcoded strings:**

```go
// ✅ CORRECT - Use generated constants
conditions := condition.New(
    condition.Condition{
        Field:    users.Id,
        Operator: condition.Equal,
        Value:    req.Id,
    },
)

// ❌ WRONG - Don't use hardcoded strings
conditions := condition.New(
    condition.Condition{
        Field:    "id",  // Hardcoded string
        Operator: condition.Equal,
        Value:    req.Id,
    },
)
```

**Benefits:**
- Type-safe - IDE can validate and autocomplete
- Refactor-friendly - Rename-safe across the codebase
- Prevents typos - Catch errors at compile time
- Consistent naming - Follows jzero conventions

## Basic Syntax

```go
conditions := condition.New(
    condition.Condition{
        Field:    users.Id,
        Operator: condition.Equal,
        Value:    value,
    },
)
```

## Comparison Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `condition.Equal` | `=` | `{Field: users.Id, Operator: condition.Equal, Value: 123}` |
| `condition.NotEqual` | `!=` / `<>` | `{Field: users.Status, Operator: condition.NotEqual, Value: "deleted"}` |
| `condition.Greater` | `>` | `{Field: users.Age, Operator: condition.Greater, Value: 18}` |
| `condition.GreaterOrEqual` | `>=` | `{Field: users.Age, Operator: condition.GreaterOrEqual, Value: 18}` |
| `condition.Less` | `<` | `{Field: users.Age, Operator: condition.Less, Value: 65}` |
| `condition.LessOrEqual` | `<=` | `{Field: users.Age, Operator: condition.LessOrEqual, Value: 10}` |

## Pattern Matching Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `condition.Like` | `LIKE` | `{Field: users.Name, Operator: condition.Like, Value: "%john%"}` |
| `condition.In` | `IN` | `{Field: users.Id, Operator: condition.In, Value: []int64{1,2,3}}` |
| `condition.NotIn` | `NOT IN` | `{Field: users.Status, Operator: condition.NotIn, Value: []string{"deleted", "banned"}}` |
| `condition.IsNull` | `IS NULL` | `{Field: users.DeletedAt, Operator: condition.IsNull}` |
| `condition.IsNotNull` | `IS NOT NULL` | `{Field: users.Email, Operator: condition.IsNotNull}` |
| `condition.Between` | `BETWEEN` | `{Field: users.CreatedAt, Operator: condition.Between, Value: []any{start, end}}` |

## Pagination & Sorting Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `condition.Limit` | `LIMIT n` | `{Operator: condition.Limit, Value: 20}` |
| `condition.Offset` | `OFFSET n` | `{Operator: condition.Offset, Value: 0}` |
| `condition.OrderBy` | `ORDER BY` | `{Operator: condition.OrderBy, Value: []string{"id DESC", "created_at ASC"}}` |

**Note**: Pagination operators are specified as conditions, not as method parameters.

## Building Complex Conditions

### Dynamic Condition Building

```go
// Build conditions dynamically
conditions := condition.New(
    // Pagination
    condition.Condition{Operator: condition.Limit, Value: 20},
    condition.Condition{Operator: condition.Offset, Value: 0},
    condition.Condition{Operator: condition.OrderBy, Value: []string{"id DESC"}},
)

// Add filter conditions conditionally
if ageFilter > 0 {
    conditions = append(conditions, condition.Condition{
        Field:    users.Age,
        Operator: condition.Equal,
        Value:    ageFilter,
    })
}

if nameSearch != "" {
    conditions = append(conditions, condition.Condition{
        Field:    users.Name,
        Operator: condition.Like,
        Value:    "%" + nameSearch + "%",
    })
}

if len(statusList) > 0 {
    conditions = append(conditions, condition.Condition{
        Field:    users.Status,
        Operator: condition.In,
        Value:    statusList,
    })
}

// Use with any *ByCondition method
users, total, err := model.PageByCondition(ctx, nil, conditions...)
```

## Using Condition Chain

For a more fluent API when building complex conditions, use `condition.NewChain()`:

### Basic Chain Usage

```go
// Build conditions with chain API
chain := condition.NewChain().
    Equal(users.Status, "active").
    GreaterOrEqual(users.Age, 18).
    Less(users.Age, 65).
    In(users.Role, []string{"admin", "user"})

// Convert to conditions slice
conditions := chain.Build()

// Use with *ByCondition methods
usersList, err := model.FindByCondition(ctx, nil, conditions...)
```

### Chain Methods

| Method | Description | Example |
|--------|-------------|---------|
| `Equal(field, value)` | `=` | `chain.Equal(users.Id, 123)` |
| `NotEqual(field, value)` | `!=` | `chain.NotEqual(users.Status, "deleted")` |
| `Greater(field, value)` | `>` | `chain.Greater(users.Age, 18)` |
| `GreaterOrEqual(field, value)` | `>=` | `chain.GreaterOrEqual(users.Age, 18)` |
| `Less(field, value)` | `<` | `chain.Less(users.Age, 65)` |
| `LessOrEqual(field, value)` | `<=` | `chain.LessOrEqual(users.Age, 100)` |
| `Like(field, value)` | `LIKE` | `chain.Like(users.Name, "%john%")` |
| `In(field, values)` | `IN` | `chain.In(users.Id, []int64{1,2,3})` |
| `NotIn(field, values)` | `NOT IN` | `chain.NotIn(users.Status, []string{"deleted"})` |
| `Between(field, min, max)` | `BETWEEN` | `chain.Between(users.CreatedAt, start, end)` |
| `IsNull(field)` | `IS NULL` | `chain.IsNull(users.DeletedAt)` |
| `IsNotNull(field)` | `IS NOT NULL` | `chain.IsNotNull(users.Email)` |

### Conditional Building with Chain

```go
// Start with base conditions
chain := condition.NewChain().
    Equal(users.Status, "active")

// Add conditions dynamically
if minAge > 0 {
    chain = chain.GreaterOrEqual(users.Age, minAge)
}

if maxAge > 0 {
    chain = chain.LessOrEqual(users.Age, maxAge)
}

if searchQuery != "" {
    chain = chain.Like(users.Name, "%"+searchQuery+"%")
}

// Build and use
conditions := chain.Build()
usersList, total, err := model.PageByCondition(ctx, nil, conditions...)
```

## Combining Chain and New

You can combine both approaches for maximum flexibility:

```go
// Use chain for filters
filters := condition.NewChain().
    Equal(users.Status, "active").
    GreaterOrEqual(users.Age, 18).
    Build()

// Combine with pagination
conditions := condition.New(
    condition.Condition{Operator: condition.Limit, Value: 20},
    condition.Condition{Operator: condition.Offset, Value: 0},
)
conditions = append(conditions, filters...)

// Use combined conditions
usersList, total, err := model.PageByCondition(ctx, nil, conditions...)
```

## When to Use Chain vs. New

- Use `condition.New()` when you need special operators like `Limit`, `Offset`, `OrderBy`
- Use `condition.NewChain()` for cleaner syntax when building filter conditions
- Combine both approaches for complex queries with filters and pagination

## Related Documentation

- [Model Generation](./model-generation.md) - Generating models with field constants
- [CRUD Operations](./crud-operations.md) - Using conditions with CRUD methods
- [Best Practices](./best-practices.md) - Database usage guidelines
