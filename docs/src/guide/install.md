---
title: 环境准备
icon: download
order: 2
---

由于 jzero 基于 go 语言以及 go-zero 框架, 需要先安装 golang, goctl 工具

::: tip  如果觉得需要安装的工具太多可以采取使用 Docker 的方式, 工具全部集成到容器中
:::

## 安装 golang

:::tip 推荐采用 go 1.21 版本以及以上
:::

使用 gvm 工具安装 golang, 并能管理 golang 的版本.

[gvm Release](https://github.com/jaronnie/gvm/releases)

```shell
# 以 linux 为例子, 下载 tar.gz 后
tar -zxvf gvm_1.4.2_Linux_x86_64.tar.gz
mv gvm /usr/local/bin
gvm init
# 重新开一个 terminal 或者 source 一下配置文件. 如 source ~/.bashrc
gvm install go1.22.2
gvm activate go1.22.2
```

## 安装 jzero

```shell
go install github.com/jzero-io/jzero@latest

jzero version
```

## jzero 相关工具一键安装

```shell
jzero check
```

## Docker

::: code-tabs#shell

@tab Docker
```shell
docker pull jaronnie/jzero:latest
```
:::