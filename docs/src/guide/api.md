---
title: 保姆级 api 教程
icon: eos-icons:api
star: true
order: 0.2
---

## 概述

api 是 go-zero 自研的领域特性语言（下文称 api 语言 或 api 描述语言）, 旨在实现人性化的基础描述语言, 作为生成 HTTP 服务最基本的描述语言.

jzero 扩展了 api 语法, 支持了一下特性: 

* `go_package`: 将 go types 生成在定义的 package 中, 所以能支持不同 api 文件的 type 定义可以同名, 保持和 proto 中的 `go_package` 一致
* `compact_handler`: 能将同一组路由的 handler 生成在同一个文件中, 减少文件的数量, 保持和 proto 的 server 模块一致

## api 定义

```api
info (
    // 定义 go_package: 生成的 go types 放入的文件夹位置
    go_package: "user"
)

type User {
    id int `json:"id"`
    username string `json:"username"`
}

type PageRequest {
    page int `form:"page"`
    size int `form:"size"`
    username string `form:"username,optional"` // 过滤参数, 可选
}

type PageResponse {
    total uint64 `json:"total"` // 总数
    list  []User `json:"list"`  // 分页数据
}

type UpdateRequest {
    id int `path:"id"`
    username string `json:"username"`
}

type UpdateResponse {}

@server (
    prefix:          /api/user // 路由 prefix
    group:           user      // 生成的 handler/logic 文件夹位置
    jwt:             JwtAuth   // 是否启用 jwt 验证
    middleware:      AuthX     // 该组路由的中间件
    compact_handler: true      // 是否合并该组的 handler 为同一个文件, 默认每个路由都有 handler 文件
)
service simpleapi {
    @doc "用户分页"
    @handler Page
    get /page (PageRequest) returns (PageResponse)

    @doc "更新用户"
    @handler Update
    post /update (UpdateRequest) returns (UpdateResponse)
}
```

对应的 curl 命令:

```shell
# 用户分页接口
curl -X GET "http://localhost:8080/api/user/page?page=1&size=10&username=test" \
  -H "Authorization: Bearer <your-jwt-token>"

# 更新用户接口
curl -X POST "http://localhost:8080/api/user/update/123" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"username": "new_username"}'

```

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

## 合并同一个 group 的 handler 为同一个文件

```shell {4}
@server (
	prefix:          /api/v1
	group:           system/user
	compact_handler: true
)
service simpleapi {
	@handler GetUserHandler
	get /system/user/getUser (GetUserRequest) returns (GetUserResponse)

	@handler DeleteUserHandler
	get /system/user/deleteUser (DeleteUserRequest) returns (DeleteUserResponse)
}
```