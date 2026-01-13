# CRUD Operations

## Overview

jzero automatically generates comprehensive CRUD methods for your models. **Use these generated methods for basic operations** - only write custom SQL for advanced scenarios that generated methods can't handle.

## Generated Methods Overview

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

## Basic CRUD Operations

### Insert

Use `InsertV2` when you need the auto-increment ID:

```go
func (l *Create) Create(req *types.CreateRequest) (*types.CreateResponse, error) {
    user := &users.Users{
        Name:  req.Name,
        Email: req.Email,
        Age:   int64(req.Age),
    }

    // ✅ Use InsertV2 to get auto-increment ID
    err := l.svcCtx.Model.Users.InsertV2(l.ctx, nil, user)
    if err != nil {
        l.Logger.Errorf("failed to insert user: %v", err)
        return nil, err
    }

    return &types.CreateResponse{
        Id: user.Id,
    }, nil
}
```

### Bulk Insert

Use `BulkInsert` for batch operations:

```go
func (l *Import) Import(users []*users.UsersModel) error {
    err := l.svcCtx.Model.Users.BulkInsert(l.ctx, nil, users)
    if err != nil {
        l.Logger.Errorf("failed to bulk insert users: %v", err)
        return err
    }
    return nil
}
```

### Find by Primary Key

```go
func (l *Get) Get(req *types.GetRequest) (*types.GetResponse, error) {
    user, err := l.svcCtx.Model.Users.FindOne(l.ctx, nil, req.Id)
    if err != nil {
        if errors.Is(err, users.ErrNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    return &types.GetResponse{
        Id:    user.Id,
        Name:  user.Name,
        Email: user.Email,
        Age:   int(user.Age),
    }, nil
}
```

### Find by Conditions

> **Note:** For detailed information on building conditions, see [Condition Builder](./condition-builder.md).

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

func (l *List) List(req *types.ListRequest) (*types.ListResponse, error) {
    // ✅ Build conditions with condition options
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

### Update by Conditions

> **Note:** For detailed information on building conditions, see [Condition Builder](./condition-builder.md).

**Method 1: Using map for simple updates**

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

func (l *Update) Update(req *types.UpdateRequest) error {
    // ✅ Build conditions with chain
    conditions := condition.NewChain().
        Equal(users.Id, req.Id).
        Build()

    updateData := map[string]any{
        string(users.Name): req.Name,
        string(users.Age):  req.Age,
    }

    _, err := l.svcCtx.Model.Users.UpdateFieldsByCondition(l.ctx, nil, updateData, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to update user: %v", err)
        return err
    }

    return nil
}
```

**Method 2: Using UpdateFieldChain for complex updates**

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

func (l *Update) Update(req *types.UpdateRequest) error {
    // ✅ Build conditions with chain
    conditions := condition.NewChain().
        Equal(users.Id, req.Id).
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

### Delete by Conditions

> **Note:** For detailed information on building conditions, see [Condition Builder](./condition-builder.md).

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/users"
)

func (l *Delete) Delete(ids []int64) error {
    // ✅ Build conditions with chain
    conditions := condition.NewChain().
        In(users.Id, ids).
        Build()

    _, err := l.svcCtx.Model.Users.DeleteByCondition(l.ctx, nil, conditions...)
    if err != nil {
        l.Logger.Errorf("failed to delete users: %v", err)
        return err
    }

    return nil
}
```

## Advanced Operations

Only write custom SQL for complex queries that generated methods can't handle:

### Complex Aggregation Query

> **Note:** For detailed information on building conditions, see [Condition Builder](./condition-builder.md).

```go
import (
    "github.com/jzero-io/jzero/core/stores/condition"
    "github.com/yourproject/internal/model/orders"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (l *GetSalesReport) GetSalesReport(req *types.SalesReportRequest) (*types.SalesReportResponse, error) {
    // ✅ Build date range conditions with chain
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

## Table Sharding Pattern

`WithTable` allows you to dynamically change the table name for sharding scenarios. It accepts a function that receives the original table name and returns the modified one.

```go
func (l *GetOrder) GetOrder(userId, orderId int64) (*orders.Orders, error) {
    // Calculate shard based on user_id
    shardId := userId % 10

    // ✅ Use WithTable to specify sharded table
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

## Related Documentation

- [Model Generation](./model-generation.md) - Generating models with CRUD methods
- [Condition Builder](./condition-builder.md) - Building query conditions
- [Database Connection](./database-connection.md) - Setting up database connections
- [Best Practices](./best-practices.md) - Database usage guidelines
