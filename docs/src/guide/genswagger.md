---
title: 生成 swagger 文档
icon: lightbulb
order: 5.1
---

## 说明

同时支持基于 go-zero api 文件和 protobuf 文件自动生成 swagger 文档

其中自动生成 swagger 的插件如下:

* [goctl-swagger](https://github.com/jzero-io/goctl-swagger)
* [protoc-gen-openapiv2](https://github.com/grpc-ecosystem/grpc-gateway/tree/main/protoc-gen-openapiv2)

> 由于 go-zero 官方仓库插件 goctl-swagger 不怎么维护了, 并在实际使用中存在一些 bug, 所以 fork 了一份进行维护.

jzero 框架默认增加了在线访问 swagger ui 的路由

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen swagger
```

@tab Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app jaronnie/jzero:latest gen swagger
```
:::

## 在线访问 swagger ui

ip:port/swagger