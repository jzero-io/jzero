# Database Connection

## Overview

jzero's `modelx` package provides database connection functionality with support for MySQL, PostgreSQL, and SQLite. You don't need to import database drivers manually - jzero handles this automatically.

## Configuration

Define your database configuration in `etc/etc.yaml`:

### MySQL Configuration

```yaml
sqlx:
    driverName: "mysql"
    dataSource: "root:password@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
```

### PostgreSQL Configuration

```yaml
sqlx:
    driverName: "pgx"
    dataSource: "postgres://user:password@127.0.0.1:5432/mydb?sslmode=disable"
```

### SQLite Configuration

```yaml
sqlx:
    driverName: "sqlite"
    dataSource: "mydb.db"
```

## Redis Configuration

For caching and session management, you can configure Redis in `etc/etc.yaml`:

### Basic Redis Configuration

```yaml
redis:
    host: "127.0.0.1:6379"
    type: "node"  # node or cluster
    pass: "yourpassword"
```

### Redis Cluster Configuration

```yaml
redis:
    host: "127.0.0.1:6379"
    type: "cluster"
    pass: "yourpassword"
    # For cluster, you can specify multiple nodes
    # host: '127.0.0.1:6379,127.0.0.1:6380,127.0.0.1:6381'
```

### Advanced Redis Options

```yaml
redis:
    host: "127.0.0.1:6379"
    type: "node"
    pass: "yourpassword"
    # Optional TLS configuration
    tls: false
```

## Config Structure

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

## Complete Configuration Example

**`etc/etc.yaml`:**

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

## Service Context Integration

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

## Connection with Cache

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

## Related Documentation

- [Model Generation](./model-generation.md) - Learn how to generate models
- [CRUD Operations](./crud-operations.md) - Database operations patterns
- [Best Practices](./best-practices.md) - Database usage guidelines

For complete documentation on modelx, see [modelx Documentation](https://docs.jzero.io/component/modelx).
