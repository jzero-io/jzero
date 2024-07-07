---
title: 镜像制作与推送
icon: rocket
star: true
order: 1
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
cd your_project
docker buildx build --platform linux/amd64,linux/arm64 --progress=plain -t your_project:latest . --push
```

## 编译单平台镜像

```shell
cd your_project
docker buildx build --platform linux/amd64 --progress=plain -t your_project:latest . --load
```