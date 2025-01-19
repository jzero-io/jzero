---
title: 环境准备
icon: marketeq:download-alt-4
order: 2
---

## 安装 jzero

::: code-tabs#shell

@tab 本地安装 jzero(推荐)

```bash
# 设置国内代理
go env -w GOPROXY=https://goproxy.cn,direct
go install github.com/jzero-io/jzero@latest

# 获取 jzero 版本信息
jzero version

# 自动下载所依赖的环境
jzero check
```

@tab 基于 Docker 使用 jzero

```bash
docker pull ghcr.io/jzero-io/jzero:latest
# 如果无法 pull
docker pull registry.cn-hangzhou.aliyuncs.com/ghcr.io/jaronnie/jzero:latest
docker tag registry.cn-hangzhou.aliyuncs.com/ghcr.io/jaronnie/jzero:latest ghcr.io/jzero-io/jzero:latest
```
:::

## jzero 命令补全(可选)

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
gvm completion bash | sudo tee /etc/bash_completion.d/gvm > /dev/null
# zsh
echo "autoload -U compinit; compinit" >> ~/.zshrc
jzero completion zsh > "${fpath[1]}/_jzero"
```

### windows

```shell
jzero completion powershell | Out-String | Invoke-Expression
```

