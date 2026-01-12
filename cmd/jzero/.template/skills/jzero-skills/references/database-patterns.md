# Database Patterns

## SQL Database with jzero

jzero provides enhanced database operations built on go-zero's `sqlx` package with support for multiple databases, flexible configuration, and enhanced model generation.

### Database Connection

jzero's `modelx` package provides database connection functionality with support for MySQL, PostgreSQL, and SQLite. You don't need to import database drivers manually - jzero handles this automatically.

#### Configuration

Define your database configuration in `etc/etc.yaml`:

```yaml
sqlx:
    driverName: "mysql"  # mysql, pgx, or sqlite
    dataSource: "root:123456@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
```

**MySQL Example:**
```yaml
sqlx:
    driverName: "mysql"
    dataSource: "root:password@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
```

**PostgreSQL Example:**
```yaml
sqlx:
    driverName: "pgx"
    dataSource: "postgres://user:password@127.0.0.1:5432/mydb?sslmode=disable"
```

**SQLite Example:**
```yaml
sqlx:
    driverName: "sqlite"
    dataSource: "mydb.db"
```

#### Redis Configuration

For caching and session management, you can configure Redis in `etc/etc.yaml`:

**Basic Redis Configuration:**
```yaml
redis:
    host: "127.0.0.1:6379"
    type: "node"  # node or cluster
    pass: "yourpassword"
```

**Redis Cluster Configuration:**
```yaml
redis:
    host: "127.0.0.1:6379"
    type: "cluster"
    pass: "yourpassword"
    # For cluster, you can specify multiple nodes
    # host: '127.0.0.1:6379,127.0.0.1:6380,127.0.0.1:6381'
```

**Advanced Redis Options:**
```yaml
redis:
    host: "127.0.0.1:6379"
    type: "node"
    pass: "yourpassword"
    # Optional TLS configuration
    tls: false
```

#### Config Structure

Define the config struct in `internal/config/config.go`:

```go
package config

import (
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
    Rest   RestConf
    Log    LogConf
    Sqlx   SqlxConf
    Redis  RedisConf
    // ... other configs
}

type SqlxConf struct {
    sqlx.SqlConf
}

type RedisConf struct {
    redis.RedisConf
}
```

**Complete Configuration Example (`etc/etc.yaml`):**

```yaml
rest:
    name: myapi
    host: 0.0.0.0
    port: 8000

log:
    serviceName: myapi
    encoding: plain
    level: info
    mode: console

sqlx:
    driverName: "mysql"
    dataSource: "root:123456@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"

redis:
    host: "127.0.0.1:6379"
    type: "node"
    pass: "123456"
```

#### Service Context Integration

Initialize the database connection in `internal/svc/servicecontext.go`:

```go
package svc

import (
    "github.com/jzero-io/jzero/core/configcenter"
    "github.com/jzero-io/jzero/core/stores/modelx"
    "github.com/jzero-io/jzero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "your-project/internal/config"
    "your-project/internal/model"
)

type ServiceContext struct {
    ConfigCenter configcenter.ConfigCenter[config.Config]
    SqlxConn     sqlx.SqlConn
    Model        model.Model
    RedisConn    *redis.Redis
    Cache        cache.Cache
}

func NewServiceContext(cc configcenter.ConfigCenter[config.Config]) *ServiceContext {
    svcCtx := &ServiceContext{
        ConfigCenter: cc,
    }

    // Connect to database
    svcCtx.SqlxConn = modelx.MustNewConn(cc.MustGetConfig().Sqlx.SqlConf)

    // Connect to Redis (optional, for caching)
    svcCtx.RedisConn = redis.MustNewRedis(cc.MustGetConfig().Redis.RedisConf)
    svcCtx.Cache = cache.NewRedisNode(svcCtx.RedisConn, errors.New("cache not found"))

    // Initialize models with optional cache
    svcCtx.Model = model.NewModel(svcCtx.SqlxConn,
        modelx.WithCachedConn(modelx.NewConnWithCache(svcCtx.SqlxConn, svcCtx.Cache)),
    )

    return svcCtx
}
```

#### Connection with Cache

For better performance, you can integrate Redis caching with your database connection:

```go
// Connect to Redis
redisConn := redis.MustNewRedis(cc.MustGetConfig().Redis.RedisConf)

// Create cache node
cacheNode := cache.NewRedisNode(redisConn, errors.New("cache not found"))

// Create cached connection
cachedConn := modelx.NewConnWithCache(sqlxConn, cacheNode)

// Initialize models with cache
model := model.NewModel(sqlxConn, modelx.WithCachedConn(cachedConn))
```

For complete documentation on modelx, see [modelx Documentation](https://docs.jzero.io/component/modelx).

## Basic SQL Operations Pattern

### ✅ Model Generation with jzero

jzero supports multiple ways to generate models:

#### Method 1: From Local SQL Files (Recommended)

Place SQL DDL files in `desc/sql/` directory:

```bash
desc/sql/
   ├── users.sql
   ├── orders.sql
   └── products.sql
```

```bash
# Generate all models
jzero gen

# Or generate only specified
jzero gen --desc desc/sql/users.sql
```

#### Method 2: From Remote Datasource

Configure in `.jzero.yaml`:

```yaml
gen:
  model-driver: mysql  # or pgx for postgres
  model-cache: true    # enable cache
  model-cache-table:
    - users            # cache specific tables
  # Use remote datasource
  model-datasource: true
  model-datasource-url: "root:123456@tcp(127.0.0.1:3306)/mydb"
  model-datasource-table:
    - users
    - orders
    - products
```

```bash
jzero gen
```

#### Method 3: Multiple Datasources

```yaml
gen:
  model-driver: mysql
  model-datasource: true
  # Multiple datasources with schema.table format
  model-datasource-url:
    - "root:123456@tcp(127.0.0.1:3306)/app_db"
    - "root:123456@tcp(127.0.0.1:3306)/log_db"
  model-datasource-table:
    - app_db.users
    - app_db.orders
    - log_db.operation_logs
```

jzero automatically generates `internal/model/model.go` that registers all models:

```go
// Code generated by jzero. DO NOT EDIT.
package model

import (
    "github.com/eddieowens/opts"
    "github.com/jzero-io/jzero/core/stores/modelx"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/yourproject/internal/model/users"
    "github.com/yourproject/internal/model/orders"
)

type Model struct {
    Users  users.UsersModel
    Orders orders.OrdersModel
}

func NewModel(conn sqlx.SqlConn, op ...opts.Opt[modelx.ModelOpts]) Model {
    return Model{
        Users:  users.NewUsersModel(conn, op...),
        Orders: orders.NewOrdersModel(conn, op...),
    }
}
```

### Example SQL Schema

```sql
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL UNIQUE,
  `age` int NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### Generated Model Methods

jzero generates these methods by default:

- **Insert** - Insert single row
- **InsertV2** - Insert single row and return auto-increment ID
- **BulkInsert** - Batch insert multiple rows
- **Update** - Update by primary key
- **Delete** - Delete by primary key
- **FindOne** - Find by primary key
- **FindByCondition** - Find by custom conditions
- **FindFieldsByCondition** - Find specific fields by conditions
- **FindOneByCondition** - Find one by conditions
- **FindOneFieldsByCondition** - Find one with specific fields
- **CountByCondition** - Count by conditions
- **PageByCondition** - Paginated query
- **UpdateFieldsByCondition** - Update specific fields
- **DeleteByCondition** - Delete by conditions
- **WithTable** - Specify table name (for sharding)

## Understanding the Condition Builder

The `condition` package provides a fluent, type-safe way to build database query conditions. It's used with all `*ByCondition` generated methods.

### ✅ Use Generated Field Constants

jzero automatically generates field constants for each model in `internal/model/<table>/<table>model_gen.go`:

```go
// internal/model/users/usersmodel_gen.go
package users

const (
    Id        condition.Field = "id"
    Name      condition.Field = "name"
    Email     condition.Field = "email"
    Age       condition.Field = "age"
    Phone     condition.Field = "phone"
    Status    condition.Field = "status"
    CreatedAt condition.Field = "created_at"
    UpdatedAt condition.Field = "updated_at"
)
```

**Always use these constants instead of hardcoded strings:**

```go
import "github.com/yourproject/internal/model/users"

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

**Benefits of using constants:**
- Type-safe - IDE can validate and autocomplete
- Refactor-friendly - Rename-safe across the codebase
- Prevents typos - Catch errors at compile time
- Consistent naming - Follows jzero conventions

### Import

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"  // Import model for field constants
)
```

### Basic Syntax

```go
conditions := condition.New(
    condition.Condition{
        Field:    users.Id,        // ✅ Use generated constant
        Operator: condition.Equal,
        Value:    value,
    },
)
```

### Common Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `condition.Equal` | `=` | `{Field: users.Id, Operator: condition.Equal, Value: 123}` |
| `condition.NotEqual` | `!=` / `<>` | `{Field: users.Status, Operator: condition.NotEqual, Value: "deleted"}` |
| `condition.Greater` | `>` | `{Field: users.Age, Operator: condition.Greater, Value: 18}` |
| `condition.GreaterOrEqual` | `>=` | `{Field: users.Age, Operator: condition.GreaterOrEqual, Value: 18}` |
| `condition.Less` | `<` | `{Field: users.Age, Operator: condition.Less, Value: 65}` |
| `condition.LessOrEqual` | `<=` | `{Field: users.Age, Operator: condition.LessOrEqual, Value: 10}` |
| `condition.Like` | `LIKE` | `{Field: users.Name, Operator: condition.Like, Value: "%john%"}` |
| `condition.In` | `IN` | `{Field: users.Id, Operator: condition.In, Value: []int64{1,2,3}}` |
| `condition.NotIn` | `NOT IN` | `{Field: users.Status, Operator: condition.NotIn, Value: []string{"deleted", "banned"}}` |
| `condition.IsNull` | `IS NULL` | `{Field: users.DeletedAt, Operator: condition.IsNull}` |
| `condition.IsNotNull` | `IS NOT NULL` | `{Field: users.Email, Operator: condition.IsNotNull}` |
| `condition.Between` | `BETWEEN` | `{Field: users.CreatedAt, Operator: condition.Between, Value: []any{start, end}}` |

### Special Operators for Pagination & Sorting

| Operator | Description | Example |
|----------|-------------|---------|
| `condition.Limit` | `LIMIT n` | `{Operator: condition.Limit, Value: 20}` |
| `condition.Offset` | `OFFSET n` | `{Operator: condition.Offset, Value: 0}` |
| `condition.OrderBy` | `ORDER BY` | `{Operator: condition.OrderBy, Value: []string{"id DESC", "created_at ASC"}}` |

**Note**: Pagination operators (`Limit`, `Offset`, `OrderBy`) are specified as conditions, not as method parameters.

### Building Complex Conditions

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

### Using Condition Chain

For a more fluent API when building complex conditions, use `condition.NewChain()`:

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

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

**Chain Methods:**

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

**Conditional Building with Chain:**

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

**When to Use Chain vs. New:**

- Use `condition.New()` when you need special operators like `Limit`, `Offset`, `OrderBy`
- Use `condition.NewChain()` for cleaner syntax when building filter conditions
- You can combine both approaches:

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

## CRUD Operations Pattern

jzero automatically generates comprehensive CRUD methods for your models. **Use these generated methods for basic operations** - only write custom SQL for advanced scenarios.

### Generated Methods Overview

jzero generates these methods for each model:

| Method | Description | When to Use |
|--------|-------------|-------------|
| `Insert(ctx, session, data)` | Insert single row | Insert without getting ID |
| `InsertV2(ctx, session, data)` | Insert + return ID | Need auto-increment ID |
| `BulkInsert(ctx, session, data)` | Batch insert | Insert multiple rows efficiently |
| `FindOne(ctx, session, id)` | Find by primary key | Query single record by ID |
| `FindByCondition(ctx, session, ...)` | Find by conditions | Query with custom conditions |
| `FindOneByCondition(ctx, session, ...)` | Find one by conditions | Query single record with conditions |
| `FindFieldsByCondition(ctx, session, ...)` | Find specific fields | Query selected columns only |
| `FindOneFieldsByCondition(ctx, session, ...)` | Find one with fields | Single record + selected columns |
| `CountByCondition(ctx, session, ...)` | Count by conditions | Get total count |
| `PageByCondition(ctx, session, ...)` | Paginated query | Pagination with conditions |
| `Update(ctx, session, data)` | Update by primary key | Update known record |
| `UpdateFieldsByCondition(ctx, session, data, ...)` | Update fields by conditions | Conditional update |
| `Delete(ctx, session, id)` | Delete by primary key | Delete known record |
| `DeleteByCondition(ctx, session, ...)` | Delete by conditions | Conditional delete |
| `WithTable(func(table) string).Method(...)` | Specify table name | Table sharding |

### ✅ Basic CRUD (Use Generated Methods)

#### Insert

```go
func (l *CreateUser) CreateUser(req *types.CreateUserRequest) (*types.CreateUserResponse, error) {
    user := &model.Users{
        Name:  req.Name,
        Email: req.Email,
        Age:   int64(req.Age),
    }

    // Use InsertV2 to get auto-increment ID
    userId, err := l.svcCtx.Model.Users.InsertV2(l.ctx, nil, user)
    if err != nil {
        l.Logger.Errorf("failed to insert user: %v", err)
        return nil, err
    }

    return &types.CreateUserResponse{
        Id: userId,
    }, nil
}
```

#### Bulk Insert

```go
func (l *ImportUsers) ImportUsers(users []*users.UsersModel) error {
    err := l.svcCtx.Model.Users.BulkInsert(l.ctx, nil, users)
    if err != nil {
        l.Logger.Errorf("failed to bulk insert users: %v", err)
        return err
    }
    return nil
}
```

#### Find by Primary Key

```go
func (l *GetUser) GetUser(req *types.GetUserRequest) (*types.GetUserResponse, error) {
    user, err := l.svcCtx.Model.Users.FindOne(l.ctx, nil, req.Id)
    if err != nil {
        if errors.Is(err, modelx.ErrNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    return &types.GetUserResponse{
        Id:    user.Id,
        Name:  user.Name,
        Email: user.Email,
        Age:   int(user.Age),
    }, nil
}
```

#### Find by Conditions (Using Generated Method)

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

func (l *ListUsers) ListUsers(req *types.ListUsersRequest) (*types.ListUsersResponse, error) {
    // Build conditions including pagination
    conditions := condition.New(
        condition.Condition{Operator: condition.Limit, Value: req.Size},
        condition.Condition{Operator: condition.Offset, Value: (req.Page - 1) * req.Size},
        condition.Condition{Operator: condition.OrderBy, Value: []string{"id DESC"}},
    )

    // Add filter conditions
    if req.Age > 0 {
        conditions = append(conditions, condition.Condition{
            Field: users.Age, Operator: condition.Equal, Value: req.Age,
        })
    }

    if req.Name != "" {
        conditions = append(conditions, condition.Condition{
            Field: users.Name, Operator: condition.Like, Value: "%" + req.Name + "%",
        })
    }

    // Use generated PageByCondition method
    users, total, err := l.svcCtx.Model.Users.PageByCondition(l.ctx, nil, conditions...)

    return &types.ListUsersResponse{List: users, Total: total}, err
}
```

#### Update by Conditions (Using Generated Method)

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) error {
    conditions := condition.New(
        condition.Condition{Field: users.Id, Operator: condition.Equal, Value: req.Id},
    )

    updateData := map[string]any{
        "name": req.Name,
        "age":  req.Age,
    }

    _, err := l.svcCtx.Model.Users.UpdateFieldsByCondition(l.ctx, nil, updateData, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to update user: %v", err)
        return err
    }

    return nil
}
```

#### Delete by Conditions (Using Generated Method)

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

func (l *DeleteUsersLogic) DeleteUsers(ids []int64) error {
    conditions := condition.New(
        condition.Condition{Field: users.Id, Operator: condition.In, Value: ids},
    )

    _, err := l.svcCtx.Model.Users.DeleteByCondition(l.ctx, nil, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to delete users: %v", err)
        return err
    }

    return nil
}
```

### 🚨 Advanced Operations (Write Custom SQL)

Only write custom SQL for complex queries that generated methods can't handle:

#### Complex Aggregation Query

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/orders"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (l *GetSalesReportLogic) GetSalesReport(req *types.SalesReportRequest) (*types.SalesReportResponse, error) {
    chain := condition.NewChain().
        GreaterOrEqual(orders.CreatedAt, req.StartDate).
        Less(orders.CreatedAt, req.EndDate)

    // Build aggregation query
    stmt, args := condition.BuildSelect(
        sqlbuilder.Select(
            "DATE(created_at) as date",
            "COUNT(*) as total_orders",
            "SUM(amount) as total_amount",
            "AVG(amount) as avg_amount",
        ).
            From("orders").
            GroupBy("DATE(created_at)").
            OrderBy("date").Desc(),
        chain.Build()...,
    )

    var reports []types.SalesReportItem
    err := l.svcCtx.SqlxConn.QueryRowsCtx(l.ctx, &reports, stmt, args...)
    if err != nil {
        return nil, err
    }

    return &types.SalesReportResponse{
        Reports: reports,
    }, nil
}
```

### Table Sharding Pattern

`WithTable` allows you to dynamically change the table name for sharding scenarios. It accepts a function that receives the original table name and returns the modified one.

```go
func (l *GetOrder) GetOrder(userId, orderId int64) (*orders.Orders, error) {
    // Calculate shard based on user_id
    shardId := userId % 10

    // Use WithTable to specify sharded table
    // WithTable takes a function that receives the original table name
    // and returns the modified table name
    order, err := l.svcCtx.Model.Orders.
        WithTable(func(table string) string {
            return fmt.Sprintf("orders_%d", shardId)
        }).
        FindOne(l.ctx, nil, orderId)

    if err != nil {
        return nil, err
    }

    return order, nil
}
```

## Best Practices Summary

### ✅ DO:
- Place SQL files in `desc/sql/` directory
- Use `jzero gen` to generate models from SQL files
- Use `.jzero.yaml` for generation configuration
- **Use generated field constants (e.g., `users.Id`) instead of hardcoded strings**
- Use condition builder for complex queries
- Always pass `context.Context` to database operations
- Use transactions for atomic operations
- Enable caching for read-heavy models via `model-cache: true` and `model-cache-table`
- Use `InsertV2` to get auto-increment IDs
- Use `BulkInsert` for batch operations
- Use `WithTable` for table sharding
- Support multiple databases with dynamic configuration
- Log database errors with context

### ❌ DON'T:
- Execute raw SQL without parameterization
- **Use hardcoded strings for field names (e.g., `Field: "id"` instead of `Field: users.Id`)**
- Ignore errors from database operations
- Use `_` to discard errors
- Create database connections in handlers/logic
- Keep transactions open longer than necessary
- Query in loops (use batch operations)
- Store sensitive data unencrypted
- Use `SELECT *` in production code
- Cache write-heavy data unnecessarily
- Forget to close result sets/cursors