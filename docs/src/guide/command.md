---
title: 命令行说明
icon: heroicons:command-line
order: 2.1
---

## jzero 命令一览

* check: 一键安装所有所需的工具
* new: 基于模板新建一个完整的项目
* gen: 生成服务端代码
  * gen sdk: 生成 http 客户端 sdk
  * gen zrpcclient: 生成 zrpc 客户端 sdk
  * gen swagger: 生成 swagger.json
* ivm: 接口版本管理
  * ivm init: 初始化新版本代码
  * ivm add: 新增可描述文件
    * ivm add api: 新增 api 文件
    * ivm add proto: 新增 proto 文件
* template: 模板功能
  * template init: 将模板下载到本地磁盘
  * template build: 将当前项目编译成 jzero 可使用的模板

### 获取更多信息

```shell
jzero -h
```

## 基于配置文件使用 jzero

:::tip
jzero version >= v0.23.0
:::

由于各个命令 flag 众多, 使用起来较为麻烦, 可通过配置文件启动 jzero 命令, 使用起来更为简洁

yaml 的编写规则是与命令行工具相匹配的, 如 gen sdk 命令下的 goModule flag 就写入 yaml gen.sdk.goModule 配置中

如果是 persistent_flag, 则需要写入到对应命令下的配置中, 如 gen sdk 命令下的 style flag, 需要定义 style 到 gen 配置中

其他命令类同

```yaml
syntax: v1

gen:
  style: go_zero
  sdk:
    goModule: github.com/jzero-io/httpsdk
    output: client-go
```

## 命令行补全

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

