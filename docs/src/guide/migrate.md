---
title: Database version automatic migration
icon: /icons/carbon-migrate.svg
star: true
order: 5.5
---

* jzero implements database migration capability based on [migrate](https://github.com/golang-migrate/migrate)
* jzero detects files under desc/sql_migration directory by default, executes migration
* Refer to [best practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) on how to write database migration files

## Configuration

```yaml
migrate:
  driver: "mysql"
  datasource-url: "root:123456@tcp(127.0.0.1:3306)/jzero-admin"
```

## Upgrade

```shell
# Upgrade to latest by default
jzero migrate up
# Upgrade n migrations
jzero migrate up 3
```

## Rollback

```shell
# Rollback all by default
jzero migrate down
# Rollback n migrations
jzero migrate down 3
```

## Get version

```shell
jzero migrate version
```

## Force rollback to specific version

```shell
jzero migrate goto <your_version>
```
