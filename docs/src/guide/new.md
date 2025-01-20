---
title: 新建项目
icon: mdi:new-box
order: 3
---

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
# 生成 swagger
jzero gen swagger
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

## jzero 官方外置模板(branch)

:::important
通过不同的分支, 可以选择不同的模板. 

[默认远程仓库地址](https://github.com/jzero-io/templates)
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

⚠️ 如果仓库需要鉴权, 请设置 remote-auth-username 和 remote-auth-password(token) 参数
:::

### 新建项目

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --remote https://gitlab.xxx.com/xx/your_repo.git --branch xxx
```

@tab Docker

```bash
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project --remote https://gitlab.xxx.com/xx/your_repo.git --branch xxx
```
:::