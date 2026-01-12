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

- **Use generated field constants (e.g., `users.Id`) instead of hardcoded strings**
  ```go
  // ✅ CORRECT
  conditions := condition.New(
      condition.Condition{Field: users.Id, Operator: condition.Equal, Value: req.Id},
  )

  // ❌ WRONG
  conditions := condition.New(
      condition.Condition{Field: "id", Operator: condition.Equal, Value: req.Id},
  )
  ```

### Query Building

- **Use condition builder for complex queries** - Provides type-safe, fluent API
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
func (l *GetUser) GetUser(req *types.GetUserRequest) (*types.GetUserResponse, error) {
    user, err := l.svcCtx.Model.Users.FindOne(l.ctx, nil, req.Id)
    if err != nil {
        if errors.Is(err, users.ErrNotFound) {
            return nil, errors.New("user not found")
        }
        // Log unexpected errors with context
        l.Logger.Errorf("failed to find user %d: %v", req.Id, err)
        return nil, err
    }
    return &types.GetUserResponse{...}, nil
}
```

### Pagination with Conditions

```go
func (l *ListUsers) ListUsers(req *types.ListUsersRequest) (*types.ListUsersResponse, error) {
    conditions := condition.New(
        condition.Condition{Operator: condition.Limit, Value: req.Size},
        condition.Condition{Operator: condition.Offset, Value: (req.Page - 1) * req.Size},
        condition.Condition{Operator: condition.OrderBy, Value: []string{"id DESC"}},
    )

    if req.Status != "" {
        conditions = append(conditions, condition.Condition{
            Field: users.Status, Operator: condition.Equal, Value: req.Status,
        })
    }

    users, total, err := l.svcCtx.Model.Users.PageByCondition(l.ctx, nil, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to list users: %v", err)
        return nil, err
    }

    return &types.ListUsersResponse{List: users, Total: total}, nil
}
```

### Batch Operations

```go
func (l *UpdateUsersStatus) UpdateUsersStatus(userIds []int64, status string) error {
    conditions := condition.New(
        condition.Condition{Field: users.Id, Operator: condition.In, Value: userIds},
    )

    updateData := map[string]any{string(users.Status): status}

    _, err := l.svcCtx.Model.Users.UpdateFieldsByCondition(l.ctx, nil, updateData, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to update users status: %v", err)
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
        if errors.Is(err, modelx.ErrNotFound) {
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
- [Condition Builder](./condition-builder.md) - Building query conditions
- [CRUD Operations](./crud-operations.md) - Using generated methods
