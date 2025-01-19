---
title: 新建项目
icon: mdi:new-box
order: 3
---

## 场景一: 基于 api 框架一键构建 http 项目

:::tip
[点击了解 api 使用教程](https://jzero.jaronnie.com/guide/develop/api.html)
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --frame api
# 生成代码
jzero gen
go mod tidy
go run main.go server
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame api
```
:::

## 场景二: 基于 zrpc 框架一键构建 rpc 项目

:::tip
[点击了解 proto 使用教程](https://jzero.jaronnie.com/guide/develop/proto.html)
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --frame rpc

# 生成代码
jzero gen
go mod tidy
go run main.go server
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame rpc --output /app/your_project
```
:::

## 场景三: 基于 gateway 框架一键构建 gateway 项目

:::tip
即提供 rpc 接口的同时, 提供 http 的方式调用 rpc 服务.
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --frame gateway

# 生成代码
jzero gen
go mod tidy
go run main.go server
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame gateway
```
:::
