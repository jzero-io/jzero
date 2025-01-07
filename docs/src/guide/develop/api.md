---
title: api 文件使用文档
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

:::tip
保证 .jzero.yaml 文件中的 gen.split-api-types-dir 配置为 true, 否则不生效
:::

```api
syntax = "v1"

info (
	go_package: "version"
)
```

jzero 脚手架推荐的 api 文件内容如下:

可以通过如下命令生成改文件:

```shell
jzero ivm add api --name user
```

```api
syntax = "v1"

info (
	go_package: "user"
)

type CreateRequest {}

type CreateResponse {}

type ListRequest {}

type ListResponse {}

type GetRequest {}

type GetResponse {}

type EditRequest {}

type EditResponse {}

type DeleteRequest {}

type DeleteResponse {}

@server (
	prefix: /api/v1
	group:  user
)
service ntls {
	@handler CreateHandler
	post /user/create (CreateRequest) returns (CreateResponse)

	@handler ListHandler
	get /user/list (ListRequest) returns (ListResponse)

	@handler GetHandler
	get /user/get (GetRequest) returns (GetResponse)

	@handler EditHandler
	post /user/edit (EditRequest) returns (EditResponse)

	@handler DeleteHandler
	get /user/delete (DeleteRequest) returns (DeleteResponse)
}
```