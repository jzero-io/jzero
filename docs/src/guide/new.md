---
title: 新建项目
icon: mdi:new-box
order: 3
---

## 说明

新建项目基于模板文件进行构建, 默认模板规范如下:

:::tip 如果构建自己的模板, app 文件夹是必须的

由于新建项目完全基于模板文件, 这意味这你可以构建任意框架的代码作为模板, 默认使用 go-zero 框架

默认远程模板仓库: [https://github.com/jzero-io/templates](https://github.com/jzero-io/templates)
:::

```shell
app: 服务端代码模板
go-zero: go-zero 框架模板
client:
  client-go: go 客户端代码模板
  client-ts: ts 客户端代码模板
docs:
  markdown: markdown 文档模板
ivm:
  init: 新版本接口初始化代码模板
  add:
    api: example api 文件模板
    proto: example proto 文件模板
```

## 新建项目命令参数

```shell
$ jzero new -h                                           
jzero new project

Usage:
  jzero new [flags]

Flags:
  -b, --branch string    remote templates repo branch
  -h, --help             help for new
      --home string      set home dir
  -m, --module string    set go module
  -o, --output string    set output dir
  -r, --remote string    remote templates repo (default "https://github.com/jzero-io/templates")
      --with-template   with template files in your project

Global Flags:
      --config string   config file (default is $HOME/.jzero/config.yaml)
      --debug           debug mode
```

| 参数字段      | 参数类型 | 是否必填 | 默认值                                | 参数说明                                    |
| ------------- | -------- | -------- | ------------------------------------- | ------------------------------------------- |
| branch        | string   | 否       | 空字符串                              | 远程仓库，配合 remote 参数使用              |
| home          | string   | 否       | 空字符串                              | 模板仓库本地路径                            |
| module        | String   | 否       | args[0]                               | go module                                   |
| output        | string   | 否       | args[0]                               | 输出文件夹路径                              |
| remote        | string   | 否       | https://github.com/jzero-io/templates | 远程仓库路                                  |
| with-template | bool     | 否       | 否                                    | 是否将模板文件放入项目中的 .template 文件夹 |


## 场景一: 基于 go-zero 的 api 框架一键构建 http 项目

:::tip 
[点击了解 go-zero api 的特性以及如何使用](https://go-zero.dev/docs/tutorials)
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch api
```

@tab Docker

```bash
docker run --rm -v ${PWD}/your_project:/app/your_project jaronnie/jzero:latest new your_project --branch api
```
:::

## 场景二: 基于 go-zero 的 zrpc 框架一键构建 rpc 项目

:::tip
[点击了解 go-zero zrpc 的特性以及如何使用](https://go-zero.dev/docs/tutorials/grpc/server/configuration)
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch rpc
# 如果需要与原生 goctl 保持同样的目录结构请使用
jzero new your_project --branch rpc-goctl
```

@tab Docker

```bash
docker run --rm -v ${PWD}/your_project:/app/your_project jaronnie/jzero:latest new your_project --branch rpc
# 如果需要与原生 goctl 保持同样的目录结构请使用
docker run --rm -v ${PWD}/your_project:/app/your_project jaronnie/jzero:latest new your_project --branch rpc-goctl
```
:::

## 场景三: 基于 go-zero 的 gateway 框架一键构建 gateway 项目

:::tip
即提供 rpc 接口的同时, 提供 http 的方式调用 rpc 服务.
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch gateway
```

@tab Docker

```bash
docker run --rm -v ${PWD}/your_project:/app/your_project jaronnie/jzero:latest new your_project --branch gateway
```
:::

## 场景四: 同时支持 go-zero api, rpc, gateway

:::tip
能通过 go-zero 的 api 框架编写路由, 同时具备 rpc 服务以及 http 调用 rpc 服务

这是最复杂的一个场景, jzero 默认模板为这个
:::

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project
```

@tab Docker

```bash
docker run --rm -v ${PWD}/your_project:/app/your_project jaronnie/jzero:latest new your_project
```
:::

## 场景五: 命令行项目

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
docker run --rm -v ${PWD}/your_project:/app/your_project jaronnie/jzero:latest new your_project --branch cli
```
:::

