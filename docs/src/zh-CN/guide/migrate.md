---
title: 数据库版本自动迁移
icon: carbon:migrate
star: true
order: 5.5
---

* jzero 基于 [migrate](https://github.com/golang-migrate/migrate) 实现数据库迁移能力
* jzero 默认检测 desc/sql_migration 目录下的文件, 执行迁移
* 参考 [最佳实践](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) 如何编写数据库迁移文件

## 配置

```yaml
migrate:
  driver: "mysql"
  datasource-url: "root:123456@tcp(127.0.0.1:3306)/jzero-admin"
```

## 升级

```shell
# 默认升级到最新
jzero migrate up
# 升级 n 个 migrations
jzero migrate up 3
```

## 回滚

```shell
# 默认回滚所有
jzero migrate down
# 回滚 n 个 migrations
jzero migrate down 3
```

## 获取版本

```shell
jzero migrate version
```

## 强制回滚到某个版本

```shell
jzero migrate goto <your_version>
```