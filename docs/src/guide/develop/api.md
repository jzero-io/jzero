---
title: api 使用文档
icon: eos-icons:api
star: true
order: 0.2
category: 开发
tag:
  - Guide
---

:::tip
[go-zero api 文档](https://go-zero.dev/docs/tutorials)
:::

## api 字段校验

> jzero 集成 [https://github.com/go-playground/validator](https://github.com/go-playground/validator) 进行字段校验

```api
syntax = "v1"

type CreateRequest {
    name string `json:"name" validate:"gte=2,lte=30"` // 名称
}
```

## api types 文件分组分文件夹

```api
syntax = "v1"

info (
	go_package: "version"
)
```

## 合并同一个 group 的 handler 为同一个文件

```api
@server (
	prefix:          /api/v1
	group:           system/user
	compact_handler: true
)
service simpleapi {
	@handler GetUserHandler
	get /system/user/getUser (GetUser2Request) returns (GetUserResponse)

	@handler DeleteUserHandler
	get /system/user/deleteUser (DeleteUserRequest) returns (DeleteUserResponse)
}
```

## 自动生成 api 文件:

```shell
jzero ivm add api --name user
```