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

## 命令补全

命令自动补全作为命令行工具的一个重要特性, 能够显著提升开发效率, 减少输入错误, 并帮助用户快速掌握工具的使用方法.

对于 jzero 这样的项目生成和开发工具来说，命令补全功能可以让用户更便捷地使用其丰富的功能, 尤其是在面对复杂的命令和参数时.

* 减少输入错误：命令补全可以避免因拼写错误而导致的命令失败
* 快速查找命令：用户可以通过部分输入快速定位到完整的命令或参数
* 降低学习成本：新用户可以通过补全提示快速了解工具支持的命令和功能

当然, jzero 为了更好的提升用户的使用体验, jzero 支持通过配置文件和环境变量的组合方式使用 jzero, 详情请查看[这里](/guide/overview.html#不同姿势使用-jzero)

对于不同操作系统与不同的 shell 类型, jzero 命令行补全的安装方式如下:

### macOS

```shell
# bash
jzero completion bash > /usr/local/etc/bash_completion.d/jzero
# zsh
echo "autoload -U compinit; compinit" >> ~/.zshrc
jzero completion zsh > "${fpath[1]}/_jzero"
```

### linux

```shell
# bash
gvm completion bash | sudo tee /etc/bash_completion.d/jzero > /dev/null
# zsh
echo "autoload -U compinit; compinit" >> ~/.zshrc
jzero completion zsh > "${fpath[1]}/_jzero"
```

### windows

```shell
jzero completion powershell | Out-String | Invoke-Expression
```

