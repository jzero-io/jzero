# Database Best Practices

## Overview

This guide outlines best practices for working with databases in jzero applications to ensure security, performance, and maintainability.

## DO ✅

### Model Generation

- **Place SQL files in `desc/sql/` directory** - Keep your DDL files organized
- **Use `jzero gen` to generate models from SQL files** - Leverage code generation
- **Use `.jzero.yaml` for generation configuration** - Centralize your generation settings
- **Enable caching for read-heavy models** - Use `model-cache: true` and `model-cache-table` for appropriate tables

### Field Constants

- **✅ Use generated field constants (e.g., `users.Id`) instead of hardcoded strings**
  ```go
  // ✅ CORRECT - Use chain API with generated constants
  conditions := condition.NewChain().
      Equal(users.Id, req.Id).
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
      Equal(users.Status, "active").
      Build()

  // ❌ WRONG - NEVER use condition.New()
  conditions := condition.New(
      condition.Condition{Field: users.Status, Operator: condition.Equal, Value: "active"},
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

```go
func (l *Get) Get(req *types.GetRequest) (*types.GetResponse, error) {
    user, err := l.svcCtx.Model.Users.FindOne(l.ctx, nil, req.Id)
    if err != nil {
        if errors.Is(err, users.ErrNotFound) {
            return nil, errors.New("user not found")
        }
        // Log unexpected errors with context
        l.Logger.Errorf("failed to find user %d: %v", req.Id, err)
        return nil, err
    }
    return &types.GetResponse{...}, nil
}
```

### Pagination with Conditions

```go
func (l *List) List(req *types.ListRequest) (*types.ListResponse, error) {
    // ✅ Build conditions with condition options
    conditions := condition.NewChain().
        Equal(users.Status, req.Status,
            condition.WithSkipFunc(func() bool {
                return req.Status == ""  // Skip if Status empty
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
        In(users.Id, userIds).
        Build()

    updateData := map[string]any{string(users.Status): status}

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
func (l *UpdateUser) UpdateUser(userId int64, req *types.UpdateUserRequest) error {
    // ✅ Build conditions with chain
    conditions := condition.NewChain().
        Equal(users.Id, userId).
        Build()

    // ✅ Build update fields with UpdateFieldChain
    updateFields := condition.NewUpdateFieldChain().
        Assign(users.Name, req.Name).
        Assign(users.Email, req.Email).
        Incr(users.Version).              // Increment version
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
func (l *GetOrder) GetOrder(userId, orderId int64) (*orders.Orders, error) {
    shardId := userId % 10

    order, err := l.svcCtx.Model.Orders.
        WithTable(func(table string) string {
            return fmt.Sprintf("orders_%d", shardId)
        }).
        FindOne(l.ctx, nil, orderId)

    if err != nil {
        if errors.Is(err, orders.ErrNotFound) {
            return nil, errors.New("order not found")
        }
        l.Logger.Errorf("failed to find order %d for user %d: %v", orderId, userId, err)
        return nil, err
    }

    return order, nil
}
```

## Related Documentation

- [Database Connection](./database-connection.md) - Setting up database connections
- [Model Generation](./model-generation.md) - Generating models with field constants
- [Condition Builder](./condition-builder.md) - Building query conditions (MUST read)
- [CRUD Operations](./crud-operations.md) - Using generated methods
