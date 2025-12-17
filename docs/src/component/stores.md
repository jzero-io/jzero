---
title: stores(存储组件)
icon: streamline-plump-color:database
order: 2
---

## migrate

```go
m, err := migrate.NewMigrate(cc.MustGetConfig().Sqlx.SqlConf)
if err != nil {
	return err
}

if err = m.Up(); err != nil {
	return err
}
```