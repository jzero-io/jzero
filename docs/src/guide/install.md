---
title: 环境准备
icon: download
order: 2
---

由于 jzero 基于 go 语言以及 go-zero 框架, 需要先安装 golang, goctl 工具

::: tip  如果觉得需要安装的工具太多可以采取使用 Docker 的方式, 工具全部集成到容器中
:::

## 安装 golang

使用 gvm 工具安装 golang, 并能管理 golang 的版本.

[gvm Release](https://github.com/jaronnie/gvm/releases)

```shell
# 以 linux 为例子, 下载 tar.gz 后
tar -zxvf gvm_1.4.2_Linux_x86_64.tar.gz
mv gvm /usr/local/bin
gvm init
# 重新开一个 terminal 或者 source 一下配置文件. 如 source ~/.bashrc
gvm install go1.19
gvm activate go1.19
```

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