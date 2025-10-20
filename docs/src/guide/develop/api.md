---
title: api 使用文档
icon: eos-icons:api
star: true
order: 0.2
category: 开发
tag:
  - Guide
---

## 前言

通过 go-zero 自研的 api 文件定义, 称为 api 可描述语言, 可用于

* 自动生成多语言服务端代码
* 自动生成文档(json/html/swagger)
* 自动生成多语言客户端代码

与 proto 文件有异曲同工之妙, 但是比 proto 更简单易用

[go-zero api 教程](https://go-zero.dev/docs/tutorials)

![](http://oss.jaronnie.com/image-20250120232337438.png)

## api 字段校验

> jzero 默认集成 [https://github.com/go-playground/validator](https://github.com/go-playground/validator) 进行字段校验

```shell {4}
syntax = "v1"

type CreateRequest {
    name string `json:"name" validate:"gte=2,lte=30"` // 名称
}
```

## 将 types 文件夹按照 go_package 进行分组

:::important go_package 的选项, 参考自 proto 文件, 能将 message 生成的结构体分组

在 api 文件中同理, go_package 选项能将定义的 type 生成的结构体分组

两大优点: 
1. 避免默认生成的 types/types.go 爆炸

2. 提升开发体验, 不同 group 下的 type 命名不会冲突
:::

```shell {3,4,5,6}
syntax = "v1"

info (
	go_package: "version"
)
```

## 合并同一个 group 的 handler/logic 为同一个文件

```shell {4,5}
@server (
	prefix:          /api/v1
	group:           system/user
	compact_handler: true
	compact_logic: true
)
service simpleapi {
	@handler GetUserHandler
	get /system/user/getUser (GetUser2Request) returns (GetUserResponse)

	@handler DeleteUserHandler
	get /system/user/deleteUser (DeleteUserRequest) returns (DeleteUserResponse)
}
```