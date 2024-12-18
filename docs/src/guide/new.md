---
title: 新建项目
icon: mdi:new-box
order: 3
---

## 说明

新建项目基于模板文件进行构建, 默认模板规范如下:

:::tip 如果构建自己的模板, app 文件夹是必须的

由于新建项目完全基于模板文件, 这意味这你可以构建任意框架的代码作为模板, 默认使用 go-zero 框架

远程模板仓库: [https://github.com/jzero-io/templates](https://github.com/jzero-io/templates)

从 v0.26.0 开始, api/gateway/rpc 模板内置, 使用 --frame 参数可以切换模板框架
:::

## 场景一: 基于 go-zero 的 api 框架一键构建 http 项目(默认 frame)

:::tip 
[点击了解 go-zero api 的特性以及如何使用](https://go-zero.dev/docs/tutorials)

[点击了解在 jzero 脚手架中 api 使用教程](https://jzero.jaronnie.com/guide/develop/api.html)
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch api # jzero version < v0.26.0
jzero new your_project # jzero version >= v0.26.0
# 如果需要与原生 goctl 保持同样的目录结构请使用
jzero new your_project --branch api-goctl
# 如果需要将应用部署在 vercel 上
jzero new your_project --branch api-vercel

# 生成代码
jzero gen
go mod tidy
go run main.go server
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --branch api
# 如果需要与原生 goctl 保持同样的目录结构请使用
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --branch api-goctl
```
:::

## 场景二: 基于 go-zero 的 zrpc 框架一键构建 rpc 项目

:::tip
[点击了解 go-zero zrpc 的特性以及如何使用](https://go-zero.dev/docs/tutorials/grpc/server/configuration)

[点击了解在 jzero 脚手架中 proto 使用教程](https://jzero.jaronnie.com/guide/develop/proto.html)
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch rpc # jzero version < v0.26.0
jzero new your_project --frame rpc # jzero version >= v0.26.0
# 如果需要与原生 goctl 保持同样的目录结构请使用
jzero new your_project --branch rpc-goctl

# 生成代码
jzero gen
go mod tidy
go run main.go server
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame rpc --output /app/your_project
# 如果需要与原生 goctl 保持同样的目录结构请使用
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --branch rpc-goctl
```
:::

## 场景三: 基于 go-zero 的 gateway 框架一键构建 gateway 项目

:::tip
即提供 rpc 接口的同时, 提供 http 的方式调用 rpc 服务.
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch gateway # jzero version < v0.26.0
jzero new your_project --frame gateway # jzero version >= v0.26.0

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

## 场景四: 命令行项目

:::tip
基于 [cobra](https://github.com/spf13/cobra) 构建命令行项目
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch cli
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --branch cli
```
:::

