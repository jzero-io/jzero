---
title: 新建项目工程
icon: mdi:new-box
order: 3
---

## 模板介绍

模板是一个预先定义好的一组代码结构，为你提供了项目的基础架构和工程规范

模板可以帮助你快速开始创建一个项目工程, 而不需要从零开始编写代码

## 特性

jzero 是一个强大的项目创建工具，支持多种场景下的项目创建需求.

它提供了多种模板和灵活的配置方式，能够满足从个人开发者到企业团队的各种需求

* 内置模板(frame): jzero 默认的模板, 仅提供框架核心能力
* 本地模板(home): 可自行编辑模板内容, 满足你的定制需求
* 远程仓库模板(branch): 可用来构建企业内部的模板仓库, 适配不同公司内部的开发需要

## jzero 内置模板(frame)

:::important 通过 go embed 特性, 内置在 jzero 命令行工具中的模板

[模板地址](https://github.com/jzero-io/jzero/tree/main/.template/frame)
:::

### 新建 api 项目

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --frame api
cd your_project
# 生成代码
jzero gen
# 生成 swagger
jzero gen swagger
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame api
cd your_project
# 生成代码
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# 生成 swagger
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# 下载依赖
go mod tidy
# 启动项目
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```
:::

### 新建 rpc 项目

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --frame rpc
# 生成代码
jzero gen
# 下载依赖
go mod tidy
# 启动项目
go run main.go server
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame rpc
cd your_project
# 生成代码
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# 下载依赖
go mod tidy
# 启动项目
go run main.go server
```
:::

### 新建 gateway 项目

:::important 即提供 rpc 接口的同时, 提供 http 的方式调用 rpc 服务.
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --frame gateway
# 生成代码
jzero gen
# 生成 swagger
jzero gen swagger
# 下载依赖
go mod tidy
# 启动项目
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --frame gateway
cd your_project
# 生成代码
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# 生成 swagger
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# 下载依赖
go mod tidy
# 启动项目
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```
:::

## 本地磁盘路径模板

:::important 指定本地磁盘路径作为模板

可使用 `jzero template init` 命令将内置模板持久化到本地磁盘

路径为: **$HOME/.jzero/templates/$VERSION**
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --home your_template_path
```
:::


## jzero 官方外置模板(branch)

:::important 通过不同的分支, 可以选择不同的模板. 

⚠️需要科学上网并在命令行模式下设置代理, 否则可能会失败.

[默认远程仓库地址](https://github.com/jzero-io/templates)
:::


:::tip 对于外置模板, 默认会从远程仓库拉取模板

将 clone 的仓库缓存到本地磁盘: **$HOME/.jzero/templates/remote/xxx**

使用 cache 参数可以不再重复拉取仓库, 如 **jzero new your_project --branch xx --cache**
:::

### 新建 go-zero 原汁原味的 api 项目

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch api-goctl
cd your_project
# 生成代码
jzero gen
# 下载依赖
go mod tidy
# 启动项目
go run main.go
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --branch api-goctl
cd your_project
# 生成代码
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# 下载依赖
go mod tidy
# 启动项目
go run main.go
```
:::

### 新建 go-zero 原汁原味的 rpc 项目

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch rpc-goctl
cd your_project
# 生成代码
jzero gen
# 下载依赖
go mod tidy
# 启动项目
go run main.go
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --branch rpc-goctl
cd your_project
# 生成代码
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# 下载依赖
go mod tidy
# 启动项目
go run main.go
```
:::

### 新建命令行项目

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch cli
cd your_project
# 下载依赖
go mod tidy
# 启动项目
go run main.go
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --branch cli
cd your_project
# 下载依赖
go mod tidy
# 启动项目
go run main.go
```
:::

## 基于企业内部模板仓库新建项目

:::important 设置 remote 仓库地址和 branch 分支参数

⚠️ 如果仓库需要鉴权, 请设置 remote-auth-username 和 remote-auth-password flag

HTTP 协议建议设置环境变量 `JZERO_REMOTE_AUTH_USERNAME` 和 `JZERO_REMOTE_AUTH_PASSWORD`

避免暴露敏感信息
:::

### 新建项目

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --remote https://gitlab.xxx.com/xx/your_repo.git --branch xxx
# 如需要权限, 使用 http 协议
jzero new your_project --remote https://gitlab.xxx.com/xx/your_repo.git --branch xxx --remote-auth-username xxx --remote-auth-password xxx

# 使用 ssh 协议
jzero new your_project --remote git@gitlab.xxx.com:xx/your_repo.git --branch xxx
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --remote https://gitlab.xxx.com/xx/your_repo.git --branch xxx
# 如需要权限, 使用 http 协议
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --remote https://gitlab.xxx.com/xx/your_repo.git --branch xxx --remote-auth-username xxx --remote-auth-password xxx
# 使用 ssh 协议
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --remote git@gitlab.xxx.com:xx/your_repo.git --branch xxx
```
:::