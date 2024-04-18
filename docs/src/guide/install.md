---
title: 环境准备
icon: download
order: 2
---

由于 jzero 基于 go-zero 框架, 需要先安装 goctl 工具

::: tip  如果觉得需要安装的工具太多可以采取使用 Docker 的方式, 工具全部集成到容器中
:::

## 安装 goctl

```shell
go install github.com/zeromicro/go-zero/tools/goctl@latest

goctl --version
```

## 安装 proto 相关工具

```shell
goctl env check --install --verbose --force
```

## 安装 jzero

```shell
go install github.com/jaronnie/jzero@latest

jzero version
```

## 安装 goreleaser

::: tip  jzero version >= v0.7.3 引入 .goreleaser.yaml
:::

```shell
go install github.com/goreleaser/goreleaser@latest
```

## 安装 task

::: tip  jzero version >= v0.7.3 引入 Taskfile.yml
:::

```shell
go install github.com/go-task/task/v3/cmd/task@latest
```

## Docker

::: code-tabs#shell

@tab Docker(amd64)
```shell
docker pull jaronnie/jzero:latest
```

@tab Docker(arm64)
```shell
docker pull jaronnie/jzero:latest-arm64
```
:::