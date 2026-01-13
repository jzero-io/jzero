# Condition Builder

## Overview

The `condition` package provides a fluent, type-safe way to build database query conditions using the **chain API**.

## Chain Methods

### Comparison Operators

| Method                             | Description | Example                                   |
|------------------------------------|-------------|-------------------------------------------|
| `Equal(field, value)`              | `=` | `chain.Equal(usersmodel.Id, 123)`              |
| `NotEqual(field, value)`           | `!=` / `<>` | `chain.NotEqual(usersmodel.Status, "deleted")` |
| `GreaterThan(field, value)`        | `>` | `chain.GreaterThan(usersmodel.Age, 18)`        |
| `GreaterThanOrEqual(field, value)` | `>=` | `chain.GreaterThanOrEqual(usersmodel.Age, 18)` |
| `LessThan(field, value)`           | `<` | `chain.LessThan(usersmodel.Age, 65)`           |
| `LessThanOrEqual(field, value)`    | `<=` | `chain.LessThanOrEqual(usersmodel.Age, 10)`    |

### Pattern Matching Operators

| Method | Description | Example |
|--------|-------------|---------|
| `Like(field, value)` | `LIKE` | `chain.Like(usersmodel.Name, "%john%")` |
| `In(field, values)` | `IN` | `chain.In(usersmodel.Id, []int64{1,2,3})` |
| `NotIn(field, values)` | `NOT IN` | `chain.NotIn(usersmodel.Status, []string{"deleted", "banned"})` |
| `IsNull(field)` | `IS NULL` | `chain.IsNull(usersmodel.DeletedAt)` |
| `IsNotNull(field)` | `IS NOT NULL` | `chain.IsNotNull(usersmodel.Email)` |
| `Between(field, min, max)` | `BETWEEN` | `chain.Between(usersmodel.CreatedAt, start, end)` |

### Pagination & Sorting

| Method | Description | Example |
|--------|-------------|---------|
| `Page(page, size)` | `LIMIT/OFFSET` | `chain.Page(1, 20)` |
| `Limit(n)` | `LIMIT n` | `chain.Limit(20)` |
| `Offset(n)` | `OFFSET n` | `chain.Offset(0)` |
| `OrderBy(fields ...string)` | `ORDER BY` | `chain.OrderBy("id DESC", "created_at ASC")` |

## Building Complex Conditions

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    usersmodel "github.com/yourproject/internal/model/users"
)

func (l *SearchUsers) SearchUsers(req *types.SearchRequest) error {
    // ✅ Build all conditions with condition options
    conditions := condition.NewChain().
        Equal(usersmodel.Status, req.Status,
            condition.WithSkipFunc(func() bool {
                return req.Status == ""  // Skip if Status empty
            }),
        ).
        GreaterThanOrEqual(usersmodel.Age, req.MinAge,
            condition.WithSkipFunc(func() bool {
                return req.MinAge <= 0  // Skip if MinAge not set
            }),
        ).
        LessThanOrEqual(usersmodel.Age, req.MaxAge,
            condition.WithSkipFunc(func() bool {
                return req.MaxAge <= 0  // Skip if MaxAge not set
            }),
        ).
        Like(usersmodel.Name, "%"+req.Name+"%",
            condition.WithSkipFunc(func() bool {
                return req.Name == ""  // Skip if Name empty
            }),
        ).
        IsNotNull(usersmodel.EmailVerifiedAt,
            condition.WithSkipFunc(func() bool {
                return !req.EmailVerified  // Skip if not verified
            }),
        ).
        In(usersmodel.Role, req.Roles,
            condition.WithSkipFunc(func() bool {
                return len(req.Roles) == 0  // Skip if Roles empty
            }),
        ).
        GreaterThanOrEqual(usersmodel.CreatedAt, req.StartDate,
            condition.WithSkipFunc(func() bool {
                return req.StartDate.IsZero()  // Skip if StartDate not set
            }),
        ).
        LessThan(usersmodel.CreatedAt, req.EndDate,
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
