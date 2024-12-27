---
title: 编译二进制文件
icon: fa6-solid:square-binary
star: true
order: 1
category: 开发
tag:
  - Guide
---

## 正常编译

```shell
GOOS=linux GOARCH=amd64 go build
```

## 优化二进制体积

```shell
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
```

## 继续优化

> 如果未使用 Kubernetes 的服务发现，可以在编译的时候使用 -tags no_k8s 来排除 k8s 相关的依赖包。

```shell
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags no_k8s
```

## 压缩二进制

[下载 upx](https://github.com/upx/upx/releases)

```shell
upx your_binary
```

实测从 go build 90MB 到最终使用 upx 压缩后, 二进制大小为 12MB.