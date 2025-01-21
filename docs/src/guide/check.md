---
title: 环境准备
icon: marketeq:download-alt-4
order: 2
---

## 安装 jzero

::: important go version >= 1.22.10
:::

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

## 不同姿势使用 jzero

:::important 涨知识的小技巧
:::

* 支持通过配置文件 .jzero.yaml 控制各种参数(**强烈推荐在每个项目的根目录新建该文件**)
* 支持通过 flag 控制各种参数
* 支持通过环境变量控制各种参数
* 支持通过以上组合的方式控制各种参数, 优先级从高到低为 环境变量  > flag  > 配置文件

如: `jzero gen --style go_zero` 对应 .jzero.yaml 内容

```yaml
gen:
  style: go_zero
```

即 `jzero gen` + `.jzero.yaml` = `jzero gen --style go_zero`

对于环境变量的使用, 需要增加前缀 `JZERO_`, 如 `JZERO_GEN_STYLE`

即 `JZERO_GEN_STYLE=go_zero jzero gen` = `jzero gen --style go_zero`

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

