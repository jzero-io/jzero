---
title: condition(数据库条件查询)
icon: material-symbols:conditions
order: 4
---

condition 的核心在于构造出 statement 和 args 参数, 然后搭配 go-zero 的 sqlx 执行器实际执行.

## 特性

* 依赖于 [go-sqlbuilder](https://github.com/huandu/go-sqlbuilder) 一套代码兼容多种常用的数据库类型
* 支持链式调用方便使用

:::tip 搭配 jzero 的数据库代码自动生成的功能, 仅需构造出 conditions 即可
:::

## 查询场景

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
	// 加载配置
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// 连接 mysql 并返回 flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	conditions := condition.New(condition.Condition{
		// 操作的字段
		Field: "name",
		// 操作
		Operator: condition.Equal,
		// 字段的值
		Value: "jzero",
		// ValueFunc 优先级比 Skip 高
		ValueFunc: func() any {
			return "jzero"
		},
		// 是否跳过该条件
		Skip: false,
		// SkipFunc 优先级比 Skip 高
		SkipFunc: func() bool {
			return false
		},
	})

	// 设置全局 flavor(默认 mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildSelect(sqlbuilder.Select("*").From("user"), conditions...)
	fmt.Println(statement, args)

	// 使用特定 flavor
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
	// 加载配置
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// 连接 mysql 并返回 flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	chain := condition.NewChain().Equal("name", "jzero",
		// WithValueFunc 比 value 优先级高
		condition.WithValueFunc(func() any {
			return "jzero"
		}),
		// 是否跳过该条件
		condition.WithSkip(false),
		// WithSkipFunc 优先级比 WithSkip 高
		condition.WithSkipFunc(
			func() bool {
				return false
			}),
	)

	// 设置全局 flavor(默认 mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildSelect(sqlbuilder.Select("*").From("user"), chain.Build()...)
	fmt.Println(statement, args)

	// 使用特定 flavor
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

## 更新场景

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
	// 加载配置
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// 连接 mysql 并返回 flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	conditions := condition.New(condition.Condition{
		// 操作的字段
		Field: "name",
		// 操作
		Operator: condition.Equal,
		// 字段的值
		Value: "jzero",
		// ValueFunc 优先级比 Skip 高
		ValueFunc: func() any {
			return "jzero"
		},
		// 是否跳过该条件
		Skip: false,
		// SkipFunc 优先级比 Skip 高
		SkipFunc: func() bool {
			return false
		},
	})

	// 设置全局 flavor(默认 mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildUpdate(
		sqlbuilder.Update("user"),
		// 设置更新字段, 直接使用 map
		map[string]any{
			"name": "jzero",
			"version": condition.UpdateField{
				Operator: condition.Incr,
			},
		},
		conditions...)

	// 使用特定 flavor
	statement, args = condition.BuildUpdateWithFlavor(flavor,
		sqlbuilder.Update("user"),
		// 设置更新字段, 直接使用 map
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
	// 加载配置
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// 连接 mysql 并返回 flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	chain := condition.NewChain().Equal("name", "jzero",
		// WithValueFunc 比 value 优先级高
		condition.WithValueFunc(func() any {
			return "jzero"
		}),
		// 是否跳过该条件
		condition.WithSkip(false),
		// WithSkipFunc 优先级比 WithSkip 高
		condition.WithSkipFunc(
			func() bool {
				return false
			}),
	)

	// 设置全局 flavor(默认 mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildUpdate(
		sqlbuilder.Update("user"),
		// 设置更新字段, 构造出 map
		condition.NewUpdateFieldChain().
			Assign("name", "jzero").
			Incr("version").
			Build(),
		chain.Build()...)
	fmt.Println(statement, args)

	// 使用特定 flavor
	statement, args = condition.BuildUpdateWithFlavor(flavor,
		sqlbuilder.Update("user"),
		// 设置更新字段, 构造出 map
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

## 删除场景

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
	// 加载配置
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// 连接 mysql 并返回 flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	conditions := condition.New(condition.Condition{
		// 操作的字段
		Field: "name",
		// 操作
		Operator: condition.Equal,
		// 字段的值
		Value: "jzero",
		// ValueFunc 优先级比 Skip 高
		ValueFunc: func() any {
			return "jzero"
		},
		// 是否跳过该条件
		Skip: false,
		// SkipFunc 优先级比 Skip 高
		SkipFunc: func() bool {
			return false
		},
	})

	// 设置全局 flavor(默认 mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildDelete(sqlbuilder.DeleteFrom("user"), conditions...)
	fmt.Println(statement, args)

	// 使用特定 flavor
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
	// 加载配置
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml"),
	)

	// 连接 mysql 并返回 flavor
	sqlConn, flavor := modelx.MustNewConnAndSqlbuilderFlavor(cc.MustGetConfig().Sqlx)

	chain := condition.NewChain().Equal("name", "jzero",
		// WithValueFunc 比 value 优先级高
		condition.WithValueFunc(func() any {
			return "jzero"
		}),
		// 是否跳过该条件
		condition.WithSkip(false),
		// WithSkipFunc 优先级比 WithSkip 高
		condition.WithSkipFunc(
			func() bool {
				return false
			}),
	)

	// 设置全局 flavor(默认 mysql)
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	statement, args := condition.BuildDelete(sqlbuilder.DeleteFrom("user"), chain.Build()...)
	fmt.Println(statement, args)

	// 使用特定 flavor
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