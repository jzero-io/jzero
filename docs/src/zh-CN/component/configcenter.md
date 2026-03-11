---
title: configcenter(配置中心)
icon: catppuccin:astro-config
order: 1
---

## fsnotify 实现

::: code-tabs#shell

@tab main.go

```go
package main

import (
	"fmt"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// 支持环境变量
	cc = configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml", subscriber.WithUseEnv(true)),
	)

	// 设置配置变更回调
	cc.AddListener(func() {})

	// 获取配置
	cfg, err := cc.GetConfig()
	if err != nil {
		panic(err)
	}

	// 必须获取配置
	cfg = cc.MustGetConfig()

	fmt.Println(cfg)
}

```

@tab etc/etc.yaml

```yaml
sqlx:
    datasource: "jzero-admin.db"
    driverName: "sqlite"
```

:::

### 使用环境变量

:::tip
参考 [envsubst](https://github.com/a8m/envsubst) 查看更多环境变量的设置方法
:::

```yaml
sqlx:
    # 从 DATASOURCE 获取 sqlx 的 datasource 配置, 未配置则为 jzero-admin.db
    datasource: "${DATASOURCE:-jzero-admin.db}"
    # 从 DRIVER_NAME 获取 sqlx 的 driverName 配置, 未配置则为 sqlite
    driverName: "${DRIVER_NAME:-sqlite}"
```

## etcd 实现

::: code-tabs#shell

@tab main.go

```go
package main

import (
	"fmt"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewEtcdSubscriber(subscriber.EtcdConf{
			Hosts: []string{"127.0.0.1:2379"},
			Key:   "jzero-admin",
		}),
	)

	// 设置配置变更回调
	cc.AddListener(func() {})

	// 获取配置
	cfg, err := cc.GetConfig()
	if err != nil {
		panic(err)
	}

	// 必须获取配置
	cfg = cc.MustGetConfig()

	fmt.Println(cfg)
}

```

:::