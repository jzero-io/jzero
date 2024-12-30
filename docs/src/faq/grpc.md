---
title: gateway/rpc 框架问题记录
icon: logos:grpc
star: true
order: 3
category: faq
tag:
  - faq
---

## 服务端 panic, 客户端会收到详细的错误信息暴露了服务端

![](http://oss.jaronnie.com/image-20241226200214351.png)

![](http://oss.jaronnie.com/image-20241226200236983.png)

解决方案:

* 去掉 grpc 内置的 recover interceptor, 改为自定义的 recover interceptor

![](http://oss.jaronnie.com/image-20241227114544243.png)

![](http://oss.jaronnie.com/image-20241227114611375.png)

代码如下:

```go
package middleware

import (
 "context"
 "runtime/debug"

 "github.com/zeromicro/go-zero/core/logx"
 "google.golang.org/grpc"
 "google.golang.org/grpc/codes"
 "google.golang.org/grpc/status"
)

// StreamRecoverInterceptor catches panics in processing stream requests and recovers.
func StreamRecoverInterceptor(svr any, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
 handler grpc.StreamHandler) (err error) {
 defer handleCrash(func(r any) {
  err = toPanicError(r)
 })

 return handler(svr, stream)
}

// UnaryRecoverInterceptor catches panics in processing unary requests and recovers.
func UnaryRecoverInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo,
 handler grpc.UnaryHandler) (resp any, err error) {
 defer handleCrash(func(r any) {
  err = toPanicError(r)
 })

 return handler(ctx, req)
}

func handleCrash(handler func(any)) {
 if r := recover(); r != nil {
  handler(r)
 }
}

func toPanicError(r any) error {
 logx.Errorf("%+v\n\n%s", r, debug.Stack())
 return status.Errorf(codes.Internal, "Service temporarily unavailable")
}
```
