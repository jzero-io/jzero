---
title: api 教程
icon: eos-icons:api
star: true
order: 0.2
category: 开发
tag:
  - Guide
---

:::tip 基于 go-zero api 框架: https://go-zero.dev/docs/tutorials
:::

## api 字段校验

> jzero 集成 [https://github.com/go-playground/validator](https://github.com/go-playground/validator) 进行字段校验

```api
syntax = "v1"

type CreateRequest {
    name string `json:"name" validate:"gte=2,lte=30"` // 名称
}
```
