---
title: 镜像推送
icon: puzzle-piece
star: true
order: 3
category: 开发
tag:
  - Guide
---

## 创建一个可以构建多平台的 buildx

```shell
docker buildx create --use --name=mybuilder --driver docker-container --driver-opt image=dockerpracticesig/buildkit:master
```

## 推送多平台镜像

```shell
cd app1
docker buildx build --platform linux/amd64,linux/arm64 --progress=plain -t app1:latest . --push
```

## 编译单平台镜像

```shell
cd app1
docker buildx build --platform linux/amd64 --progress=plain -t app1:latest . --load
```