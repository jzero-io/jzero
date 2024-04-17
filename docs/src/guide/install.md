---
title: 环境准备
icon: download
order: 2
---

由于 jzero 基于 go-zero 框架, 需要先安装 goctl 工具

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