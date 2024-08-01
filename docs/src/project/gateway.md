---
title: gateway 项目实战
icon: mdi:arrow-projectile-multiple
star: true
order: 1
category: project
tag:
  - Guide
---

## 1. 简介

gateway 项目包含以下技术栈:
* cobra: 提供命令行框架
* zrpc: 基于 go-zero zrpc 框架提供 rpc 服务
* gateway: 基于 go-zero gateway 框架, 支持以 http 方式调用 rpc 服务

## 2. 新建项目

```shell
jzero new simplegateway --branch gateway
cd simplegateway
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

![](https://oss.jaronnie.com/image-20240731134511973.png)

至此, 你已经拥有了一个可用的 gateway 服务.

## 5. 开发教程

### 5.1. 鉴权

在实际的项目中, 鉴权是必不可少的, 那么在 jzero 的 gateway 模板中, 如何进行鉴权呢?

gateway 模板内置了中间件 `header_processor`

位置在 `internal/middlewares/header_processor.go`

内容如下:

```go
package middleware

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/gateway"
	"google.golang.org/grpc"
)

func WithHeaderProcessor() gateway.Option {
	return gateway.WithHeaderProcessor(func(header http.Header) []string {
		var headers []string
		//// You can add header from request header here
		//// for example
		//for k, v := range header {
		//	if k == "Authorization" {
		//		headers = append(headers, fmt.Sprintf("%s:%s", k, strings.Join(v, ";")))
		//	}
		//}
		return headers
	})
}

func WithUnaryInterceptorValue(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	//md, b := metadata.FromIncomingContext(ctx)
	//if !b {
	//	return handler(ctx, req)
	//}
	//// You can verify Authorization here and set user info in context value
	//// get Authorization
	//value := md.Get("Authorization")
	//if len(value) == 1 {
	//	// set context value
	//	ctx = context.WithValue(ctx, "Authorization", value[0])
	//}
	return handler(ctx, req)
}
```

在 `WithHeaderProcessor` 中, 可将前端传来的 `Authorization` header 放入 grpc 的 metadata 中. 修改为如下:

```go
func WithHeaderProcessor() gateway.Option {
	return gateway.WithHeaderProcessor(func(header http.Header) []string {
		var headers []string
		// You can add header from request header here
		// for example
		for k, v := range header {
			if k == "Authorization" {
				headers = append(headers, fmt.Sprintf("%s:%s", k, strings.Join(v, ";")))
			}
		}
		return headers
	})
}
```

在 grpc 拦截器中获取 `Authorization` 值, 并校验获取出用户信息, 并设置到 context value 中.

```go
func WithUnaryInterceptorValue(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	md, b := metadata.FromIncomingContext(ctx)
	if !b {
		return handler(ctx, req)
	}
	// You can verify Authorization here and set user info in context value
	// get Authorization
	value := md.Get("Authorization")
	if len(value) == 1 {
		// set context value
		ctx = context.WithValue(ctx, "User", "jzero")
	}
	return handler(ctx, req)
}

```

在 `internal/middleware.go` 中, 添加该拦截器

```go
func RegisterZrpc(z *zrpc.RpcServer) {
	z.AddUnaryInterceptors(ServerValidationUnaryInterceptor)
	z.AddUnaryInterceptors(WithUnaryInterceptorValue)
}
```

后续在 logic 获取用户信息如下:

```go
func (l *SayHello) SayHello(in *hellopb.SayHelloRequest) (*hellopb.SayHelloResponse, error) {
	user := l.ctx.Value("User")
	
	return &hellopb.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s", user),
	}, nil
}
```
