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
chain := condition.NewChain().
    Equal(users.Id, value)

// Convert to conditions slice
conditions := chain.Build()

// Use with any *ByCondition method
usersList, err := model.FindByCondition(ctx, nil, conditions...)
```

## Chain Methods

### Comparison Operators

| Method | Description | Example |
|--------|-------------|---------|
| `Equal(field, value)` | `=` | `chain.Equal(users.Id, 123)` |
| `NotEqual(field, value)` | `!=` / `<>` | `chain.NotEqual(users.Status, "deleted")` |
| `Greater(field, value)` | `>` | `chain.Greater(users.Age, 18)` |
| `GreaterOrEqual(field, value)` | `>=` | `chain.GreaterOrEqual(users.Age, 18)` |
| `Less(field, value)` | `<` | `chain.Less(users.Age, 65)` |
| `LessOrEqual(field, value)` | `<=` | `chain.LessOrEqual(users.Age, 10)` |

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
| `Limit(n)` | `LIMIT n` | `chain.Limit(20)` |
| `Offset(n)` | `OFFSET n` | `chain.Offset(0)` |
| `OrderBy(fields ...string)` | `ORDER BY` | `chain.OrderBy("id DESC", "created_at ASC")` |

## Building Complex Conditions

### Basic Chain Usage

```go
// ✅ Build multiple conditions with chain API
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

### Dynamic Condition Building

```go
// ✅ Start with base conditions
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

if len(statusList) > 0 {
    chain = chain.In(users.Status, statusList)
}

// Build and use with pagination
conditions := chain.Limit(20).Offset(0).OrderBy("id DESC").Build()
usersList, total, err := model.PageByCondition(ctx, nil, conditions...)
```

### Pagination Example

```go
func (l *List) List(req *types.ListRequest) (*types.ListResponse, error) {
    // ✅ Build conditions with pagination using chain
    chain := condition.NewChain()

    // Add filter conditions dynamically
    if req.Age > 0 {
        chain = chain.Equal(users.Age, req.Age)
    }

    if req.Name != "" {
        chain = chain.Like(users.Name, "%"+req.Name+"%")
    }

    // Add pagination and ordering
    conditions := chain.
        Limit(req.Size).
        Offset((req.Page - 1) * req.Size).
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
    // ✅ Build all conditions with chain
    chain := condition.NewChain()

    // Status filter
    if req.Status != "" {
        chain = chain.Equal(users.Status, req.Status)
    }

    // Age range
    if req.MinAge > 0 {
        chain = chain.GreaterOrEqual(users.Age, req.MinAge)
    }
    if req.MaxAge > 0 {
        chain = chain.LessOrEqual(users.Age, req.MaxAge)
    }

    // Name search
    if req.Name != "" {
        chain = chain.Like(users.Name, "%"+req.Name+"%")
    }

    // Email verification
    if req.EmailVerified {
        chain = chain.IsNotNull(users.EmailVerifiedAt)
    }

    // Role filtering
    if len(req.Roles) > 0 {
        chain = chain.In(users.Role, req.Roles)
    }

    // Created date range
    if !req.StartDate.IsZero() {
        chain = chain.GreaterOrEqual(users.CreatedAt, req.StartDate)
    }
    if !req.EndDate.IsZero() {
        chain = chain.Less(users.CreatedAt, req.EndDate)
    }

    // Add pagination and sort
    conditions := chain.
        Limit(req.Size).
        Offset((req.Page - 1) * req.Size).
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
