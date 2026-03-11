---
title: cache(缓存连接)
icon: octicon:cache-16
order: 2.1
---

## 特性

* 接口设计, 后期好扩展
* 支持 cache prefix
* 默认使用 redis 实现 cache 接口

::: code-tabs#shell

@tab main.go

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/cache"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	Redis redis.RedisConf
}

func main() {
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// 连接 redis
	rds, err := redis.NewRedis(cc.MustGetConfig().Redis)
	if err != nil {
		panic(err)
	}

	// 从 redis 新建缓存, 设置默认缓存时间 5 秒
	redisCache := cache.NewRedisNode(rds, errors.New("cache not found"), cache.WithExpiry(time.Duration(5)*time.Second))
	
	// 带 cachePrefix
	redisCache = cache.NewRedisNodeWithCachePrefix(rds, errors.New("cache not found"), "jzero:",cache.WithExpiry(time.Duration(5)*time.Second))

	// 获取 key 为 name 的数据

	var value string
	if err = redisCache.GetCtx(context.Background(), "name", &value); err != nil {
		panic(err)
	}

	fmt.Println(value)
}

```

@tab etc/etc.yaml
```yaml
redis:
    host: "127.0.0.1:6379"
    type: "node"
```

:::

