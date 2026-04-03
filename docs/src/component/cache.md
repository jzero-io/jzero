---
title: cache(Cache connection)
icon: /icons/octicon-cache-16.svg
order: 2.1
---

## Features

* Interface design, easy to extend later
* Supports cache prefix
* Uses redis implementation by default for cache interface

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

	// connect redis
	rds, err := redis.NewRedis(cc.MustGetConfig().Redis)
	if err != nil {
		panic(err)
	}

	// create cache from redis, set default cache time 5 seconds
	redisCache := cache.NewRedisNode(rds, errors.New("cache not found"), cache.WithExpiry(time.Duration(5)*time.Second))

	// with cachePrefix
	redisCache = cache.NewRedisNodeWithCachePrefix(rds, errors.New("cache not found"), "jzero:",cache.WithExpiry(time.Duration(5)*time.Second))

	// get data with key name

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
