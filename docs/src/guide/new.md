---
title: 新建项目
icon: clone
order: 3
---

## new command flag

```shell
$ jzero new -h                                           
jzero new project

Usage:
  jzero new [flags]

Flags:
      --app-dir string   set app dir (default ".")
  -b, --branch string    remote templates repo branch
  -h, --help             help for new
      --home string      set home dir
  -m, --module string    set go module
  -o, --output string    set output dir
  -r, --remote string    remote templates repo (default "https://github.com/jzero-io/templates")

Global Flags:
      --config string   config file (default is $HOME/.jzero/config.yaml)
      --debug           debug mode
```

| 参数字段 | 参数类型 | 是否必填 | 默认值                                | 参数说明                       |
| -------- | -------- | -------- | ------------------------------------- | ------------------------------ |
| app-dir  | string   | 否       | .                                     | 输出文件夹路径: output/app-dir |
| branch   | string   | 否       | 空字符串                              | 远程仓库，配合 remote 参数使用 |
| home     | string   | 否       | 空字符串                              | 模板仓库本地路径               |
| module   | String   | 否       | args[0]                               | go module                      |
| output   | string   | 否       | args[0]                               | 输出文件夹路径                 |
| remote   | string   | 否       | https://github.com/jzero-io/templates | 远程仓库路                     |

## go-zero with grpc + gateway + api

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

## go-zero with grpc + gateway

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

## go-zero with api

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

## go-zero with zrpc

::: code-tabs#shell

@tab jzero

```bash
jzero new your_project --branch rpc
```

@tab Docker

```bash
docker run --rm -v ${PWD}/your_project:/app/your_project jaronnie/jzero:latest new your_project --branch rpc
```
:::

## with cobra cli project

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

