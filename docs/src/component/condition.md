---
title: condition(Database condition query)
icon: /icons/material-symbols-conditions.svg
order: 4
---

The core of condition is to construct statement and args parameters, then use go-zero's sqlx executor to actually execute.

## Features

* Depends on [go-sqlbuilder](https://github.com/huandu/go-sqlbuilder) for one codebase compatible with multiple common database types
* Supports chain calls for easy use

:::tip Pair with jzero's automatic database code generation feature, only need to construct conditions
:::

## Query scenarios

::: code-tabs#shell

@tab condition

```go
package main

import (
	"context"
	"fmt"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero/core/stores/condition"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	// load config
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// connect mysql and return flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	conditions := condition.New(condition.Condition{
		// field to operate
		Field: "name",
		// operation
		Operator: condition.Equal,
		// field value
		Value: "jzero",
		// ValueFunc has higher priority than Skip
		ValueFunc: func() any {
			return "jzero"
		},
		// whether to skip this condition
		Skip: false,
		// SkipFunc has higher priority than Skip
		SkipFunc: func() bool {
			return false
		},
	})

	// set global flavor(default mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildSelect(sqlbuilder.Select("*").From("user"), conditions...)
	fmt.Println(statement, args)

	// use specific flavor
	statement, args = condition.BuildSelectWithFlavor(flavor, sqlbuilder.Select("id", "name").From("user"), conditions...)

	type User struct {
		Id   int64  `db:"id"`
		Name string `db:"name"`
	}

	var users []User

	err := sqlConn.QueryRowsCtx(context.Background(), &users, statement, args)
	if err != nil {
		panic(err)
	}

	fmt.Println(users)
}

```

@tab condition chain

```go
package main

import (
	"context"
	"fmt"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero/core/stores/condition"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	// load config
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// connect mysql and return flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	chain := condition.NewChain().Equal("name", "jzero",
		// WithValueFunc has higher priority than value
		condition.WithValueFunc(func() any {
			return "jzero"
		}),
		// whether to skip this condition
		condition.WithSkip(false),
		// WithSkipFunc has higher priority than WithSkip
		condition.WithSkipFunc(
			func() bool {
				return false
			}),
	)

	// set global flavor(default mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildSelect(sqlbuilder.Select("*").From("user"), chain.Build()...)
	fmt.Println(statement, args)

	// use specific flavor
	statement, args = condition.BuildSelectWithFlavor(flavor, sqlbuilder.Select("id", "name").From("user"), chain.Build()...)

	type User struct {
		Id   int64  `db:"id"`
		Name string `db:"name"`
	}

	var users []User

	err := sqlConn.QueryRowsCtx(context.Background(), &users, statement, args)
	if err != nil {
		panic(err)
	}

	fmt.Println(users)
}
```

@tab etc/etc.yaml

```yaml
sqlx:
    datasource: "jzero-admin.db"
    driverName: "sqlite"
```

:::

## Update scenarios

::: code-tabs#shell

@tab condition

```go
package main

import (
	"context"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/condition"
	"github.com/jzero-io/jzero/core/stores/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	// load config
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// connect mysql and return flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	conditions := condition.New(condition.Condition{
		// field to operate
		Field: "name",
		// operation
		Operator: condition.Equal,
		// field value
		Value: "jzero",
		// ValueFunc has higher priority than Skip
		ValueFunc: func() any {
			return "jzero"
		},
		// whether to skip this condition
		Skip: false,
		// SkipFunc has higher priority than Skip
		SkipFunc: func() bool {
			return false
		},
	})

	// set global flavor(default mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildUpdate(
		sqlbuilder.Update("user"),
		// set update fields, directly use map
		map[string]any{
			"name": "jzero",
			"version": condition.UpdateField{
				Operator: condition.Incr,
			},
		},
		conditions...)

	// use specific flavor
	statement, args = condition.BuildUpdateWithFlavor(flavor,
		sqlbuilder.Update("user"),
		// set update fields, directly use map
		map[string]any{
			"name": "jzero",
			"version": condition.UpdateField{
				Operator: condition.Incr,
			},
		},
		conditions...)

	result, err := sqlConn.ExecCtx(context.Background(), statement, args)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
```

@tab condition chain

```go
package main

import (
	"context"
	"fmt"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero/core/stores/condition"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	// load config
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// connect mysql and return flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	chain := condition.NewChain().Equal("name", "jzero",
		// WithValueFunc has higher priority than value
		condition.WithValueFunc(func() any {
			return "jzero"
		}),
		// whether to skip this condition
		condition.WithSkip(false),
		// WithSkipFunc has higher priority than WithSkip
		condition.WithSkipFunc(
			func() bool {
				return false
			}),
	)

	// set global flavor(default mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildUpdate(
		sqlbuilder.Update("user"),
		// set update fields, construct map
		condition.NewUpdateFieldChain().
			Assign("name", "jzero").
			Incr("version").
			Build(),
		chain.Build()...)
	fmt.Println(statement, args)

	// use specific flavor
	statement, args = condition.BuildUpdateWithFlavor(flavor,
		sqlbuilder.Update("user"),
		// set update fields, construct map
		condition.NewUpdateFieldChain().
			Assign("name", "jzero").
			Incr("version").
			Build(),
		chain.Build()...)

	result, err := sqlConn.ExecCtx(context.Background(), statement, args)
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

## Delete scenarios

::: code-tabs#shell

@tab condition

```go
package main

import (
	"context"
	"fmt"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero/core/stores/condition"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	// load config
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// connect mysql and return flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	conditions := condition.New(condition.Condition{
		// field to operate
		Field: "name",
		// operation
		Operator: condition.Equal,
		// field value
		Value: "jzero",
		// ValueFunc has higher priority than Skip
		ValueFunc: func() any {
			return "jzero"
		},
		// whether to skip this condition
		Skip: false,
		// SkipFunc has higher priority than Skip
		SkipFunc: func() bool {
			return false
		},
	})

	// set global flavor(default mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildDelete(sqlbuilder.DeleteFrom("user"), conditions...)
	fmt.Println(statement, args)

	// use specific flavor
	statement, args = condition.BuildDeleteWithFlavor(flavor, sqlbuilder.DeleteFrom("user"), conditions...)

	result, err := sqlConn.ExecCtx(context.Background(), statement, args)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

```

@tab condition chain

```go
package main

import (
	"context"
	"fmt"

	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/modelx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero/core/stores/condition"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	// load config
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// connect mysql and return flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	chain := condition.NewChain().Equal("name", "jzero",
		// WithValueFunc has higher priority than value
		condition.WithValueFunc(func() any {
			return "jzero"
		}),
		// whether to skip this condition
		condition.WithSkip(false),
		// WithSkipFunc has higher priority than WithSkip
		condition.WithSkipFunc(
			func() bool {
				return false
			}),
	)

	// set global flavor(default mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildDelete(sqlbuilder.DeleteFrom("user"), chain.Build()...)
	fmt.Println(statement, args)

	// use specific flavor
	statement, args = condition.BuildDeleteWithFlavor(flavor, sqlbuilder.DeleteFrom("user"), chain.Build()...)

	result, err := sqlConn.ExecCtx(context.Background(), statement, args)
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
