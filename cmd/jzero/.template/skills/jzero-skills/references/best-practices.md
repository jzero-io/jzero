# Database Best Practices

## Overview

This guide outlines best practices for working with databases in jzero applications to ensure security, performance, and maintainability.

## ⚠️ Critical Rules

These rules are critical and MUST be followed to avoid common bugs and security issues.

### 1. Model Import Rule

**‼️ ALL `internal/model/xx` imports MUST use alias `xxmodel`**

This prevents naming conflicts and makes code more maintainable.

#### ❌ WRONG - Direct import without alias
```go
import "github.com/yourproject/internal/model/users"

// ❌ Don't use users.Id
conditions := condition.NewChain().
    Equal(users.Id, req.Id).
    Build()
```

#### ✅ CORRECT - Import with alias
```go
import usersmodel "github.com/yourproject/internal/model/users"

// ✅ Use usersmodel.Id
conditions := condition.NewChain().
    Equal(usersmodel.Id, req.Id).
    Build()
```

**This applies to ALL model imports:** `usersmodel`, `ordersmodel`, `productmodel`, etc.

---

### 2. Error Handling Rule

**‼️ Always use `errors.Is()` to check errors, never use `==` comparison**

**Requirements:**
- Import `errors` from `github.com/pkg/errors` (NOT the standard library)
- Use `errors.Is(err, xxxmodel.ErrNotFound)` for error comparison

#### ❌ WRONG - Direct equality comparison
```go
import "errors"  // ❌ WRONG - standard library

// ❌ WRONG - Direct comparison
if err == usermodel.ErrNotFound {
    return nil, errors.New("用户不存在")
}
```

#### ✅ CORRECT - Use errors.Is with pkg/errors
```go
import "github.com/pkg/errors"  // ✅ CORRECT

// ✅ CORRECT - Use errors.Is()
if errors.Is(err, usermodel.ErrNotFound) {
    return nil, errors.New("用户不存在")
}
```

---

### 3. Update Method Rule

**‼️ `Update()` method requires FULL object update - partial updates NOT supported**

The generated `Update()` method performs a complete update of all fields in the model. It does NOT support optional field updates.

#### ❌ WRONG - Trying to update only some fields with Update()
```go
// ❌ WRONG - This will update ALL fields including zero values!
user := &usersmodel.Users{
    Id:   req.Id,
    Name: req.Name,
    // Email and Age will be set to zero values!
}
err := l.svcCtx.Model.Users.Update(l.ctx, nil, user)
```

#### ✅ CORRECT - Use Update() when you have the complete object
```go
// ✅ CORRECT - Update when you have the full object (e.g., after modification)
user, err := l.svcCtx.Model.Users.FindOne(l.ctx, nil, req.Id)
if err != nil {
    return err
}

// Modify fields
user.Name = req.NewName
user.Age = req.NewAge

// Update the entire object
err = l.svcCtx.Model.Users.Update(l.ctx, nil, user)
```

#### ✅ CORRECT - Use UpdateFieldsByCondition for partial updates
```go
// ✅ CORRECT - Update specific fields only
conditions := condition.NewChain().
    Equal(usersmodel.Id, req.Id).
    Build()

updateData := map[string]any{
    string(usersmodel.Name): req.Name,
    // Only Name field will be updated
}

_, err := l.svcCtx.Model.Users.UpdateFieldsByCondition(l.ctx, nil, updateData, conditions...)
```

**Summary:**
- `Update(ctx, session, data)` - Full object update (ALL fields including zero values)
- `UpdateFieldsByCondition(ctx, session, data, conditions...)` - Partial field update (only specified fields)

---

### 4. Condition Builder Rule

**‼️ ALWAYS use `condition.NewChain()` instead of `condition.New()`**

The chain API provides a fluent, type-safe interface for building conditions.

#### ❌ WRONG - Using condition.New()
```go
// ❌ WRONG - Verbose and error-prone
conditions := condition.New(
    condition.Condition{Field: usersmodel.Status, Operator: condition.Equal, Value: "active"},
)
```

#### ✅ CORRECT - Using condition.NewChain()
```go
// ✅ CORRECT - Clean and fluent API
conditions := condition.NewChain().
    Equal(usersmodel.Status, "active").
    Build()
```

---

### 5. Field Constants Rule

**‼️ ALWAYS use generated field constants (e.g., `usersmodel.Id`) instead of hardcoded strings**

Generated constants provide type safety and prevent typos.

#### ❌ WRONG - Hardcoded strings
```go
// ❌ WRONG - Hardcoded string (typo-prone)
conditions := condition.NewChain().
    Equal("id", req.Id).
    Equal("name", req.Name).
    Build()
```

#### ✅ CORRECT - Generated constants
```go
// ✅ CORRECT - Type-safe constants
conditions := condition.NewChain().
    Equal(usersmodel.Id, req.Id).
    Equal(usersmodel.Name, req.Name).
    Build()
```

---

### 6. FindOne Result Rule

**‼️ `FindOne`/`FindOneByXx` methods only need `err` check - no need to check for `nil`**

When `err == nil`, the result is guaranteed to be valid.

#### ❌ WRONG - Unnecessary nil check
```go
user, err := l.svcCtx.Model.Users.FindOne(l.ctx, nil, req.Id)
if err != nil {
    return nil, err
}
// ❌ WRONG - Unnecessary check
if user == nil {
    return nil, errors.New("user not found")
}
```

#### ✅ CORRECT - Only check err
```go
user, err := l.svcCtx.Model.Users.FindOne(l.ctx, nil, req.Id)
if err != nil {
    if errors.Is(err, usersmodel.ErrNotFound) {
        return nil, errors.New("user not found")
    }
    return nil, err
}
// ✅ CORRECT - user is valid when err == nil, no nil check needed
return &types.GetResponse{...}, nil
```

---

## Additional Best Practices

### Model Generation

- **Place SQL files in `desc/sql/` directory** - Keep your DDL files organized
- **Use `jzero gen` to generate models from SQL files** - Leverage code generation
- **Use `.jzero.yaml` for generation configuration** - Centralize your generation settings
- **Enable caching for read-heavy models** - Use `model-cache: true` and `model-cache-table` for appropriate tables

### Performance

- **Use `InsertV2` to get auto-increment IDs** - Avoids additional query
- **Use `BulkInsert` for batch operations** - More efficient than individual inserts
- **Query selected columns only** - Avoid `SELECT *` in production
- **Use batch operations instead of querying in loops** - Reduce database round-trips
- **Minimize transaction scope** - Keep transactions open only as long as necessary
- **Consider cache invalidation costs** - Don't cache write-heavy data unnecessarily

### Security

- **Always use parameterized queries** - Never execute raw SQL without parameterization
- **Encrypt sensitive data** - Protect passwords, tokens, etc.
- **Validate and sanitize user input** - Never trust user input directly

### Code Quality

- **Never ignore errors** - Always handle errors properly
- **Use service context for database connections** - Don't create connections in handlers/logic
- **Log database errors with context** - Include relevant information for debugging
- **Always pass `context.Context` to database operations** - Enables cancellation and timeout

### Database Support

- **Support multiple databases with dynamic configuration** - Make your app flexible
- **Use `WithTable` for table sharding** - Scale your data horizontally

---

## Related Documentation

- [Database Connection](./database-connection.md) - Setting up database connections
- [Model Generation](./model-generation.md) - Generating models with field constants
- [Condition Builder](./condition-builder.md) - Building query conditions (MUST read)
- [CRUD Operations](./crud-operations.md) - Using generated methods
