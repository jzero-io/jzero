---
title: api 项目实战
icon: mdi:arrow-projectile-multiple
star: true
order: 1
category: project
tag:
  - Guide
---

## 1. 简介

api 项目包含以下技术栈:
* cobra: 提供命令行框架
* api: 基于 go-zero api 框架提供 http 服务

## 2. 新建项目

```shell
jzero new simpleapi --frame api
cd simpleapi
jzero gen
go mod tidy
```

## 3. 生成 swagger

```shell
jzero gen swagger
```

## 4. 测试

```shell
go run main.go server
# 访问 localhost:8001/swagger 进行测试
```

至此, 你已经拥有了一个可用的 api 服务.

## 5. 开发教程

### 5.1. 鉴权

在实际的项目中, 鉴权是必不可少的, 那么在 jzero 的 api 模板中, 如何进行鉴权呢?

例如增加了一个鉴权中间件 AuthMiddleware, 获取鉴权信息, 然后能在 logic 中使用

```go
func(next http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        // 获取用户信息
        rctx := r.Context()
        rctx = context.WithValue(rctx, "auth", "xx")
        // 携带 auth
        next(w, r.WithContext(rctx))
	}
}
```