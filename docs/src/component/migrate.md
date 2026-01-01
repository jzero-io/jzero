---
title: migrate(管理数据库迁移脚本)
icon: streamline-plump-color:database
order: 3
---

migrate 组件默认读取 desc/sql_migration 下的 sql 文件, 管理 sql.

* Up: 默认升级所有 up 脚本, 支持传入 steps 参数升级几个
* Down: 默认回滚所有 down 脚本, 支持传入 steps 参数回滚几个
* Goto: 切到某个版本
* Version: 获取当前版本

::: code-tabs#shell

@tab main.go

```go
package main

import (
	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/migrate"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml", subscriber.WithUseEnv(true)),
	)

	m, err := migrate.NewMigrate(cc.MustGetConfig().Sqlx)
	if err != nil {
		panic(err)
	}
	
	defer m.Close()

	if err = m.Up(); err != nil {
		panic(err)
	}
}
```

@tab etc/etc.yaml

```yaml
sqlx:
    datasource: "jzero-admin.db"
    driverName: "sqlite"
```

@tab desc/sql_migration/1_init.up.sql
```sql
DROP TABLE IF EXISTS `manage_user`;

CREATE TABLE `manage_user` (
                               `id` bigint NOT NULL AUTO_INCREMENT,
                               `uuid` varchar(36) NOT NULL UNIQUE,
                               `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                               `username` varchar(30) NOT NULL,
                               `password` varchar(100) NOT NULL,
                               `nickname` varchar(30) NOT NULL,
                               `gender` varchar(1) NOT NULL,
                               `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
                               `status` varchar(1) NOT NULL,
                               `email` varchar(100) NOT NULL,
                               PRIMARY KEY (`id`),
                               UNIQUE KEY `uni_manage_user_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `manage_user` (`uuid`, `create_time`, `update_time`, `username`, `password`, `nickname`, `gender`, `phone`, `status`, `email`)
VALUES
    ('1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f','2024-10-24 09:45:00','2024-10-31 09:40:13','admin','123456','超级管理员','1','','1','');
```

@tab desc/sql_migration/1_init.down.sql

```sql
DROP TABLE IF EXISTS `manage_user`;
```
:::

## 支持多数据库切换

:::tip 根据 driver 区分不同的数据库的迁移目录
:::

migrate 增加可选参数 `WithSourceAppendDriver`:

* mysql source: desc/sql_migration/mysql
* pgx source: desc/sql_migration/pgx
* sqlite source: desc/sql_migration/sqlite

::: code-tabs#shell

@tab main.go

```go
package main

import (
	"github.com/jzero-io/jzero/core/configcenter"
	"github.com/jzero-io/jzero/core/configcenter/subscriber"
	"github.com/jzero-io/jzero/core/stores/migrate"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Config struct {
	Sqlx sqlx.SqlConf
}

func main() {
	cc := configcenter.MustNewConfigCenter[Config](
		configcenter.Config{Type: "yaml"},
		subscriber.MustNewFsnotifySubscriber("etc/etc.yaml", subscriber.WithUseEnv(true)),
	)

	m, err := migrate.NewMigrate(cc.MustGetConfig().Sqlx, migrate.WithSourceAppendDriver(true))
	if err != nil {
		panic(err)
	}
	
	defer m.Close()

	if err = m.Up(); err != nil {
		panic(err)
	}
}

```

@tab etc/etc.yaml

```yaml
sqlx:
    datasource: "jzero-admin.db"
    driverName: "sqlite"
```

@tab desc/sql_migration/sqlite/1_init.up.sql
```sql
DROP TABLE IF EXISTS `manage_user`;

CREATE TABLE `manage_user` (
                               `id` bigint NOT NULL AUTO_INCREMENT,
                               `uuid` varchar(36) NOT NULL UNIQUE,
                               `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                               `username` varchar(30) NOT NULL,
                               `password` varchar(100) NOT NULL,
                               `nickname` varchar(30) NOT NULL,
                               `gender` varchar(1) NOT NULL,
                               `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
                               `status` varchar(1) NOT NULL,
                               `email` varchar(100) NOT NULL,
                               PRIMARY KEY (`id`),
                               UNIQUE KEY `uni_manage_user_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `manage_user` (`uuid`, `create_time`, `update_time`, `username`, `password`, `nickname`, `gender`, `phone`, `status`, `email`)
VALUES
    ('1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f','2024-10-24 09:45:00','2024-10-31 09:40:13','admin','123456','超级管理员','1','','1','');
```

@tab desc/sql_migration/sqlite/1_init.down.sql

```sql
DROP TABLE IF EXISTS `manage_user`;
```
:::