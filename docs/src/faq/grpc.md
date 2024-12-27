---
title: gateway/rpc 框架问题记录
icon: vscode-icons:file-type-swagger
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