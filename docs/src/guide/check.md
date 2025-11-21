---
title: 环境准备
icon: marketeq:download-alt-4
order: 2
---

## 安装 golang

jzero 依赖于 golang 环境, 推荐使用 go1.23 版本

推荐采用 [gvm](https://github.com/jaronnie/gvm) 安装 golang 环境

## 安装 jzero

提供以下三种方式使用 jzero, 请根据实际情况任选一种即可

* 源码安装(**go version >= go1.24.3**)
* 直接[下载 jzero 二进制文件](https://github.com/jzero-io/jzero/releases)
* 基于 Docker 安装 jzero, [镜像地址](https://github.com/jzero-io/jzero/pkgs/container/jzero)

### 源码安装 jzero

```bash
# 设置国内代理
go env -w GOPROXY=https://goproxy.cn,direct
go install github.com/jzero-io/jzero/cmd/jzero@latest

# 获取 jzero 版本信息
jzero version

# 自动下载所依赖的环境
jzero check
```

### 直接下载 jzero 二进制文件

[下载地址](https://github.com/jzero-io/jzero/releases)

根据自己的操作系统选择对应的压缩包, 解压后放在 `$GOPATH/bin` 目录下即可

执行以下内容完成 jzero 的环境准备

```shell
# 获取 jzero 版本信息
jzero version

# 自动下载所依赖的环境
jzero check
```

### 基于 Docker 使用 jzero

好处便是将所有依赖的工具链全部集成在容器中, 减少用户环境依赖, 减少用户环境配置的复杂度

**github 镜像源**

```shell
# 获取 jzero 版本信息
docker run --rm ghcr.io/jzero-io/jzero:latest version
```