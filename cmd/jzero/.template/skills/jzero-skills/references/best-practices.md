# Database Best Practices

## Overview

This guide outlines best practices for working with databases in jzero applications to ensure security, performance, and maintainability.

## ⚠️ Critical Import Rule

**‼️ ALL `internal/model/xx` imports MUST use alias `xxmodel`**

### ❌ WRONG - Direct import without alias
```go
import "github.com/yourproject/internal/model/users"

// ❌ Don't use users.Id
conditions := condition.NewChain().
    Equal(users.Id, req.Id).
    Build()
```

### ✅ CORRECT - Import with alias
```go
import usersmodel "github.com/yourproject/internal/model/users"

// ✅ Use usersmodel.Id
conditions := condition.NewChain().
    Equal(usersmodel.Id, req.Id).
    Build()
```

**This applies to ALL model imports:** `usersmodel`, `ordersmodel`, `productmodel`, etc.

---

## DO ✅

### Model Generation

- **Place SQL files in `desc/sql/` directory** - Keep your DDL files organized
- **Use `jzero gen` to generate models from SQL files** - Leverage code generation
- **Use `.jzero.yaml` for generation configuration** - Centralize your generation settings
- **Enable caching for read-heavy models** - Use `model-cache: true` and `model-cache-table` for appropriate tables

### Field Constants

- **✅ Use generated field constants (e.g., `usersmodel.Id`) instead of hardcoded strings**
  ```go
  // ✅ CORRECT - Use chain API with generated constants
  conditions := condition.NewChain().
      Equal(usersmodel.Id, req.Id).
      Build()

  // ❌ WRONG - Don't use hardcoded strings
  conditions := condition.NewChain().
      Equal("id", req.Id).  // Hardcoded string
      Build()
  ```

### Query Building

- **✅ Use condition chain API for ALL query building** - Provides fluent, type-safe API
  ```go
  // ✅ CORRECT - ALWAYS use condition.NewChain()
  conditions := condition.NewChain().
      Equal(usersmodel.Status, "active").
      Build()

  // ❌ WRONG - NEVER use condition.New()
  conditions := condition.New(
      condition.Condition{Field: usersmodel.Status, Operator: condition.Equal, Value: "active"},
  )
  ```
- **Always pass `context.Context` to database operations** - Enables cancellation and timeout
- **Use transactions for atomic operations** - Maintain data consistency

### Performance

- **Use `InsertV2` to get auto-increment IDs** - Avoids additional query
- **Use `BulkInsert` for batch operations** - More efficient than individual inserts
- **Query selected columns only** - Avoid `SELECT *` in production
- **Use batch operations instead of querying in loops** - Reduce database round-trips

### Database Support

- **Support multiple databases with dynamic configuration** - Make your app flexible
- **Use `WithTable` for table sharding** - Scale your data horizontally

### Error Handling

- **Log database errors with context** - Include relevant information for debugging
- **Handle `ErrNotFound` appropriately** - Return meaningful errors to users
- **Never ignore errors** - Always check and handle database errors
- **`FindOne`/`FindOneByXx` only need `err` check** - No need to check if result is `nil` since it's guaranteed valid when `err == nil`

## DON'T ❌

### Security

- **Execute raw SQL without parameterization** - Always use parameterized queries
- **Store sensitive data unencrypted** - Encrypt passwords, tokens, etc.
- **Trust user input directly** - Always validate and sanitize

### Code Quality

- **‼️ Use `condition.New()` instead of `condition.NewChain()`** - ALWAYS use chain API
- **Use hardcoded strings for field names** - Use generated constants
- **Use `_` to discard errors** - Always handle errors
- **Create database connections in handlers/logic** - Use service context
- **Check if `FindOne`/`FindOneByXx` result is `nil`** - Only check `err`, result is valid when `err == nil`

### Performance

- **Keep transactions open longer than necessary** - Minimize transaction scope
- **Query in loops** - Use batch operations
- **Use `SELECT *` in production code** - Query only needed columns
- **Cache write-heavy data unnecessarily** - Consider cache invalidation costs

### Resource Management

- **Forget to close result sets/cursors** - Let jzero handle this with generated methods
- **Create too many connections** - Use connection pooling

## Common Patterns

### Proper Error Handling

> **Note:** `FindOne`/`FindOneByXx` methods only need to check `err`. When `err == nil`, the result is guaranteed to be valid, so no need to check for `nil`.

```go
func (l *Get) Get(req *types.GetRequest) (*types.GetResponse, error) {
    user, err := l.svcCtx.Model.Users.FindOne(l.ctx, nil, req.Id)
    if err != nil {
        if errors.Is(err, usersmodel.ErrNotFound) {
            return nil, errors.New("user not found")
        }
        // Log unexpected errors with context
        l.Logger.Errorf("failed to find user %d: %v", req.Id, err)
        return nil, err
    }
    // ✅ No nil check needed - user is valid when err == nil
    return &types.GetResponse{...}, nil
}
```

### Pagination with Conditions

```go
func (l *List) List(req *types.ListRequest) (*types.ListResponse, error) {
    // ✅ Build conditions with condition options
    conditions := condition.NewChain().
        Equal(usersmodel.Status, req.Status,
            condition.WithSkipFunc(func() bool {
                return req.Status == ""  // Skip if Status empty
            }),
        ).
        Like(usersmodel.Name, "%"+req.Name+"%",
            condition.WithSkipFunc(func() bool {
                return req.Name == ""  // Skip if Name empty
            }),
        ).
        Page(req.Page, req.Size).
        OrderBy("id DESC").
        Build()

    users, total, err := l.svcCtx.Model.Users.PageByCondition(l.ctx, nil, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to list users: %v", err)
        return nil, err
    }

    return &types.ListResponse{List: users, Total: total}, nil
}
```

### Batch Operations

```go
func (l *UpdateStatus) UpdateStatus(userIds []int64, status string) error {
    // ✅ Build conditions with chain
    conditions := condition.NewChain().
        In(usersmodel.Id, userIds).
        Build()

    updateData := map[string]any{string(usersmodel.Status): status}

    _, err := l.svcCtx.Model.Users.UpdateFieldsByCondition(l.ctx, nil, updateData, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to update users status: %v", err)
        return err
    }

    return nil
}
```

Alternatively, you can use `UpdateFieldsByCondition` with `UpdateFieldChain` for more complex update operations:

```go
func (l *Update) Update(userId int64, req *types.UpdateRequest) error {
    // ✅ Build conditions with chain
    conditions := condition.NewChain().
        Equal(usersmodel.Id, userId).
        Build()

    // ✅ Build update fields with UpdateFieldChain
    updateFields := condition.NewUpdateFieldChain().
        Assign(usersmodel.Name, req.Name).
        Assign(usersmodel.Email, req.Email).
        Incr(usersmodel.Version).              // Increment version
        Build()

    _, err := l.svcCtx.Model.Users.UpdateFieldsByCondition(l.ctx, nil, updateFields, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to update user: %v", err)
        return err
    }

    return nil
}
```

### Table Sharding

```go
import (
    ordersmodel "github.com/yourproject/internal/model/orders"
)

func (l *GetOrder) GetOrder(userId, orderId int64) (*ordersmodel.Orders, error) {
    shardId := userId % 10

    order, err := l.svcCtx.Model.Orders.
        WithTable(func(table string) string {
            return fmt.Sprintf("orders_%d", shardId)
        }).
        FindOne(l.ctx, nil, orderId)

    if err != nil {
        if errors.Is(err, ordersmodel.ErrNotFound) {
            return nil, errors.New("order not found")
        }
        l.Logger.Errorf("failed to find order %d for user %d: %v", orderId, userId, err)
        return nil, err
    }

    // ✅ No nil check needed - order is valid when err == nil
    return order, nil
}
```

## Related Documentation

- [Database Connection](./database-connection.md) - Setting up database connections
- [Model Generation](./model-generation.md) - Generating models with field constants
- [Condition Builder](./condition-builder.md) - Building query conditions (MUST read)
- [CRUD Operations](./crud-operations.md) - Using generated methods
