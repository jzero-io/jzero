---
title: mongo 使用文档
icon: logos:mongodb
star: true
order: 5
category: 开发
tag:
  - Guide
---

## 前言

jzero 支持通过指定 mongo type 将代码到 `internal/mongo` 下.

为了在使用上更加方便, jzero 自动生成了 `internal/mongo/model.go` 文件, 用于注册所有生成的 mongo model.

## 特性

* 支持 redis/sync_map 缓存

## 生成代码

```yaml
gen:
    # 指定 mongo type
    mongo-type: ["user", "role", "menu"]
    # 是否需要缓存
    mongo-cache: true
    # 指定哪些类型需要缓存
    mongo-cache-type:
      - user
```

```shell
jzero gen
```