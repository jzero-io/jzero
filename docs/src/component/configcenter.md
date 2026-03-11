---
title: configcenter(Configuration center)
icon: catppuccin:astro-config
order: 1
---

## fsnotify implementation

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

	// support environment variables
	cc = configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml", subscriber.WithUseEnv(true)),
	)

	// set config change callback
	cc.AddListener(func() {})

	// get config
	cfg, err := cc.GetConfig()
	if err != nil {
		panic(err)
	}

	// must get config
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

### Using environment variables

:::tip
Refer to [envsubst](https://github.com/a8m/envsubst) for more environment variable settings
:::

```yaml
sqlx:
    # get sqlx datasource from DATASOURCE, defaults to jzero-admin.db if not set
    datasource: "${DATASOURCE:-jzero-admin.db}"
    # get sqlx driverName from DRIVER_NAME, defaults to sqlite if not set
    driverName: "${DRIVER_NAME:-sqlite}"
```

## etcd implementation

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

	// set config change callback
	cc.AddListener(func() {})

	// get config
	cfg, err := cc.GetConfig()
	if err != nil {
		panic(err)
	}

	// must get config
	cfg = cc.MustGetConfig()

	fmt.Println(cfg)
}

```

:::
