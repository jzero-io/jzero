---
title: 项目初始化
icon: mdi:new-box
order: 3
---

## 模板介绍

模板是一个预先定义好的一组代码结构，提供了项目的基础架构和工程规范

模板可以帮助你快速开始初始化一个项目, 而不需要从零开始编写代码

## 模板类型

jzero 提供了以下几种类型模板, 满足各种场景:

* 内置模板(frame): 内置模板, 提供框架核心能力, 支持可选特性(数据库/缓存)
* 路径模板(home): 指定路径作为模板, 一般放入特定项目内部, 满足特定项目需要
* 本地模板(local): 本地全局模板, 位于 ~/.jzero/templates/local 文件夹中
* 远程仓库模板(remote+branch): 可用来构建企业专属的远程模板仓库

具体使用请参阅: [模板指南](../guide/template.md)

## 初始化 api 项目

::: code-tabs#shell

@tab jzero cli

```bash
jzero new your_project --frame api
cd your_project
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```

@tab jzero Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame api
cd your_project
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```
:::

## 初始化 rpc 项目

::: code-tabs#shell

@tab jzero cli

```bash
jzero new your_project --frame rpc
cd your_project
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
```

@tab jzero Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame rpc
cd your_project
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
```
:::

## 初始化 gateway 项目

:::important 同时支持 grpc/http 接口
:::

::: code-tabs#shell

@tab jzero cli

```bash
jzero new your_project --frame gateway
cd your_project
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```

@tab jzero Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame gateway
cd your_project
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```
:::

## 可选特性 model/redis/model+redis

基于可选特性, 提供了一整套使用 model/redis/model 的解决方案

```shell
# 使用场景: 需要连接关系型数据库(model)且包含数据库缓存(cache), redis
jzero new your_project --features model,cache,redis

# 使用场景: 需要连接关系型数据库(model), redis
jzero new your_project --features model,redis

# 使用场景: 需要连接关系型数据库(model)且包含数据库缓存(cache)
jzero new your_project --features model,cache

# 使用场景: 需要连接关系型数据库(model)
jzero new your_project --features model
```