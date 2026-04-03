---
title: modelx(Database connection)
icon: /icons/oui-vis-query-sql.svg
order: 2
---

## Features

* Adapts to mysql/postgres/sqlite, no need to import drivers

::: code-tabs#shell

@tab main.go

```go
package main

import (
	"context"
	"fmt"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/modelx"
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

	sqlConn := modelx.MustNewConn(cc.MustGetConfig().Sqlx)

	// connect database
	sqlConn := modelx.MustNewConn(cc.MustGetConfig().Sqlx)

	// execute sql
	result, err := sqlConn.ExecCtx(context.Background(), "select 1")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

```

@tab etc/etc.yaml
```yaml
sqlx:
    datasource: "jzero-admin.db"
    driverName: "sqlite"
```
:::
