---
title: 生成 swagger 文档
icon: vscode-icons:file-type-swagger
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

![](https://oss.jaronnie.com/image-20240731134511973.png)

## 高级教程

### 1. 如何自定义字段

please see: [swagger 自定义字段](../faq/swagger.md)

### 2. 基于 proto 文件自定义 header

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    info: {
        version: "v1";
    };
    security_definitions: {
        security: {
            key: "ApiKeyAuth",
            value: {
                description: "JWT token for authorization ( Bearer token )"
                type: TYPE_API_KEY,
                name: "Authorization",
                in: IN_HEADER
            },
        }
    }
};
```