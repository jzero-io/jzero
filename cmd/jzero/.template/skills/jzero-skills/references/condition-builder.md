# Condition Builder

## Overview

The `condition` package provides a fluent, type-safe way to build database query conditions using the **chain API**.

‼️ **IMPORTANT: You MUST use the `condition.NewChain()` API for all query conditions. Do NOT use `condition.New()`.**

## Import

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"  // Import model for field constants
)
```

## ✅ Use Condition Chain API

**‼️ ALWAYS use `condition.NewChain()` for building conditions - this is the ONLY supported approach.**

```go
// ✅ CORRECT - Use chain API
conditions := condition.NewChain().
    Equal(users.Id, req.Id).
    Build()

// ❌ WRONG - NEVER use condition.New()
conditions := condition.New(
    condition.Condition{Field: users.Id, Operator: condition.Equal, Value: req.Id},
)
```

## Use Generated Field Constants

**‼️ ALWAYS use generated field constants instead of hardcoded strings:**

```go
// ✅ CORRECT - Use generated constants with chain
conditions := condition.NewChain().
    Equal(users.Id, req.Id).
    Build()

// ❌ WRONG - Don't use hardcoded strings
conditions := condition.NewChain().
    Equal("id", req.Id).  // Hardcoded string
    Build()
```

**Benefits:**
- Type-safe - IDE can validate and autocomplete
- Refactor-friendly - Rename-safe across the codebase
- Prevents typos - Catch errors at compile time
- Consistent naming - Follows jzero conventions

## Basic Syntax

```go
// ✅ Build conditions with chain API
conditions := condition.NewChain().
    Equal(users.Id, value).
	Build()

// Use with any *ByCondition method
usersList, err := model.FindByCondition(ctx, nil, conditions...)
```

## Chain Methods

### Comparison Operators

| Method                             | Description | Example                                   |
|------------------------------------|-------------|-------------------------------------------|
| `Equal(field, value)`              | `=` | `chain.Equal(users.Id, 123)`              |
| `NotEqual(field, value)`           | `!=` / `<>` | `chain.NotEqual(users.Status, "deleted")` |
| `GreaterThan(field, value)`        | `>` | `chain.GreaterThan(users.Age, 18)`        |
| `GreaterThanOrEqual(field, value)` | `>=` | `chain.GreaterThanOrEqual(users.Age, 18)` |
| `LessThan(field, value)`           | `<` | `chain.LessThan(users.Age, 65)`           |
| `LessThanOrEqual(field, value)`    | `<=` | `chain.LessThanOrEqual(users.Age, 10)`    |

### Pattern Matching Operators

| Method | Description | Example |
|--------|-------------|---------|
| `Like(field, value)` | `LIKE` | `chain.Like(users.Name, "%john%")` |
| `In(field, values)` | `IN` | `chain.In(users.Id, []int64{1,2,3})` |
| `NotIn(field, values)` | `NOT IN` | `chain.NotIn(users.Status, []string{"deleted", "banned"})` |
| `IsNull(field)` | `IS NULL` | `chain.IsNull(users.DeletedAt)` |
| `IsNotNull(field)` | `IS NOT NULL` | `chain.IsNotNull(users.Email)` |
| `Between(field, min, max)` | `BETWEEN` | `chain.Between(users.CreatedAt, start, end)` |

### Pagination & Sorting

| Method | Description | Example |
|--------|-------------|---------|
| `Page(page, size)` | `LIMIT/OFFSET` | `chain.Page(1, 20)` |
| `Limit(n)` | `LIMIT n` | `chain.Limit(20)` |
| `Offset(n)` | `OFFSET n` | `chain.Offset(0)` |
| `OrderBy(fields ...string)` | `ORDER BY` | `chain.OrderBy("id DESC", "created_at ASC")` |

## Building Complex Conditions

### Basic Chain Usage

```go
// ✅ Build multiple conditions with chain API
conditions := condition.NewChain().
    Equal(users.Status, "active").
    GreaterThanOrEqual(users.Age, 18).
    LessThan(users.Age, 65).
    In(users.Role, []string{"admin", "user"}).
    Build()

// Use with *ByCondition methods
usersList, err := model.FindByCondition(ctx, nil, conditions...)
```

### Dynamic Condition Building

> **✅ RECOMMENDED: Use condition options for dynamic conditions** - This provides a cleaner, more maintainable approach compared to if statements.

```go
// ✅ Build dynamic conditions with condition options
conditions := condition.NewChain().
    Equal(users.Status, "active").
    GreaterThanOrEqual(users.Age, minAge,
        condition.WithSkipFunc(func() bool {
            return minAge <= 0  // Skip if minAge not set
        }),
    ).
    LessThanOrEqual(users.Age, maxAge,
        condition.WithSkipFunc(func() bool {
            return maxAge <= 0  // Skip if maxAge not set
        }),
    ).
    Like(users.Name, "%"+searchQuery+"%",
        condition.WithSkipFunc(func() bool {
            return searchQuery == ""  // Skip if searchQuery empty
        }),
    ).
    In(users.Status, statusList,
        condition.WithSkipFunc(func() bool {
            return len(statusList) == 0  // Skip if statusList empty
        }),
    ).
    Page(1, 20).
    OrderBy("id DESC").
    Build()

// Use with *ByCondition methods
usersList, total, err := model.PageByCondition(ctx, nil, conditions...)
```

### Pagination Example

```go
func (l *List) List(req *types.ListRequest) (*types.ListResponse, error) {
    // ✅ Build conditions with pagination using condition options
    conditions := condition.NewChain().
        Equal(users.Age, req.Age,
            condition.WithSkipFunc(func() bool {
                return req.Age <= 0  // Skip if Age not set
            }),
        ).
        Like(users.Name, "%"+req.Name+"%",
            condition.WithSkipFunc(func() bool {
                return req.Name == ""  // Skip if Name empty
            }),
        ).
        Page(req.Page, req.Size).
        OrderBy("id DESC").
        Build()

    // Use generated PageByCondition method
    users, total, err := l.svcCtx.Model.Users.PageByCondition(l.ctx, nil, conditions...)

    return &types.ListResponse{List: users, Total: total}, err
}
```

## Complete Example

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

func (l *SearchUsers) SearchUsers(req *types.SearchRequest) error {
    // ✅ Build all conditions with condition options
    conditions := condition.NewChain().
        Equal(users.Status, req.Status,
            condition.WithSkipFunc(func() bool {
                return req.Status == ""  // Skip if Status empty
            }),
        ).
        GreaterThanOrEqual(users.Age, req.MinAge,
            condition.WithSkipFunc(func() bool {
                return req.MinAge <= 0  // Skip if MinAge not set
            }),
        ).
        LessThanOrEqual(users.Age, req.MaxAge,
            condition.WithSkipFunc(func() bool {
                return req.MaxAge <= 0  // Skip if MaxAge not set
            }),
        ).
        Like(users.Name, "%"+req.Name+"%",
            condition.WithSkipFunc(func() bool {
                return req.Name == ""  // Skip if Name empty
            }),
        ).
        IsNotNull(users.EmailVerifiedAt,
            condition.WithSkipFunc(func() bool {
                return !req.EmailVerified  // Skip if not verified
            }),
        ).
        In(users.Role, req.Roles,
            condition.WithSkipFunc(func() bool {
                return len(req.Roles) == 0  // Skip if Roles empty
            }),
        ).
        GreaterThanOrEqual(users.CreatedAt, req.StartDate,
            condition.WithSkipFunc(func() bool {
                return req.StartDate.IsZero()  // Skip if StartDate not set
            }),
        ).
        LessThan(users.CreatedAt, req.EndDate,
            condition.WithSkipFunc(func() bool {
                return req.EndDate.IsZero()  // Skip if EndDate not set
            }),
        ).
        Page(req.Page, req.Size).
        OrderBy("created_at DESC").
        Build()

    // Execute query
    usersList, total, err := l.svcCtx.Model.Users.PageByCondition(l.ctx, nil, conditions...)
    // ...
}
```

## Related Documentation

- [Model Generation](./model-generation.md) - Generating models with field constants
- [CRUD Operations](./crud-operations.md) - Using conditions with CRUD methods
- [Best Practices](./best-practices.md) - Database usage guidelines
