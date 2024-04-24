---
title: 编译镜像
icon: puzzle-piece
star: true
order: 3
category: 开发
tag:
  - Guide
---

采用 goreleaser 工具交叉编译二进制文件

采用 Task 工具代替 Makefile

将这两个工具结合起来使用, 能更加方便的管理项目

[安装 goreleaser](../install.md#安装-goreleaser)

[安装 task](../install.md#安装-task)

## 编译 linux/amd64 镜像

```shell
task build:amd64

docker build -t jaronnie/jzero:latest .
```

## 编译 linux/arm64 镜像

```shell
task build:arm64

docker build -t jaronnie/jzero:latest-arm64 -f Dockerfile-arm64 .
```