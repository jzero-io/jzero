---
title: 命令行说明
icon: heroicons:command-line
order: 2.1
---

## 命令

* check: 一键安装所有所需的工具
* new: 基于模板新建一个完整的项目
* gen: 生成服务端代码
* ivm: 接口版本管理
* template: 模板的新增与初始化到本地

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

