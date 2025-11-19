---
title: 生成客户端
icon: carbon:sdk
order: 5
---

## 生成 Swagger 文档

基于 API 定义生成 Swagger 文档，提供在线 API 文档界面。

### 使用方法

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen swagger
```

@tab Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
```
:::

### 在线访问 Swagger UI

生成 Swagger 文档后，可以通过以下地址访问 Swagger UI：

**Swagger UI 地址**: `localhost:8001/swagger`

![Swagger UI 示例](https://oss.jaronnie.com/image-20240731134511973.png)

## 生成 RPC 客户端

基于 Proto 文件生成 gRPC 客户端代码。

### 使用方法

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen zrpcclient
```

@tab Docker
```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen zrpcclient
```
:::

生成的客户端代码将包含完整的 gRPC 客户端实现，支持服务调用和连接管理。