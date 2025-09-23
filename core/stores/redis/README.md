# Redis DB Support

## Overview

The Redis store now supports multiple databases (DB 0-15) for all connection types:
- **Node**: Single Redis instance with configurable DB
- **Sentinel**: Redis Sentinel with configurable DB 
- **Cluster**: Redis Cluster (DB support for cache key consistency, but Redis Cluster only uses DB 0)

## Key Features

1. **Backward Compatibility**: Default DB is 0, maintaining existing behavior
2. **Connection Isolation**: Different databases use separate connection pools
3. **Optimized Cache Keys**: DB 0 doesn't include suffix for compatibility
4. **Multiple Configuration Methods**: Both configuration-based and programmatic setup

## Usage Examples

### Configuration-based Usage

```go
// Default DB 0 (backward compatible)
conf := RedisConf{
    Host: "localhost:6379",
    Type: "node",
}
redis := MustNewRedis(conf)

// Specific DB
conf := RedisConf{
    Host: "localhost:6379",
    Type: "node", 
    DB: 2,  // Use database 2
}
redis := MustNewRedis(conf)

// Sentinel with DB
conf := RedisConf{
    Host: "sentinel1:26379,sentinel2:26379",
    Type: "sentinel",
    MasterName: "mymaster",
    DB: 1,  // Use database 1
}
redis := MustNewRedis(conf)
```

### Programmatic Usage

```go
// Using WithDB option
redis := New("localhost:6379", WithDB(3))

// Combined with other options
redis := New("localhost:6379", 
    WithUser("myuser"),
    WithPass("mypass"),
    WithDB(5),
    WithTLS(),
)
```

### JSON Configuration

```json
{
    "host": "localhost:6379",
    "type": "node",
    "db": 2,
    "user": "myuser",
    "pass": "mypass"
}
```

## Cache Key Behavior

- **DB 0**: `localhost:6379` (no suffix for backward compatibility)
- **DB 1**: `localhost:6379#1`
- **DB 15**: `localhost:6379#15`

This ensures that:
1. Existing code using DB 0 continues to work without changes
2. Different databases use separate connection pools
3. Resource management is properly isolated

## Validation

- DB field accepts values 0-15 (standard Redis database range)
- Invalid DB values will be caught during configuration validation
- Each DB gets its own connection pool for proper isolation

## Notes

- **Redis Cluster**: While DB parameter is supported in configuration for consistency, Redis Cluster mode only supports DB 0
- **Sentinel**: Full DB support with proper failover handling
- **Node**: Full DB support for single Redis instances