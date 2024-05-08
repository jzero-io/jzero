---
title: jzero 命令行
icon: terminal
order: 2.1
---

## 获取 jzero usage

```shell
$ jzero -h
jzero framework.

Usage:
  jzero [command]

Available Commands:
  check       jzero env check
  completion  Generate completion script
  app      jzero app
  gen         jzero gen code
  gensdk      jzero gensdk
  help        Help about any command
  init        jzero init
  new         jzero new project
  template    jzero template
  version     jzero version

Flags:
      --config string   config file (default is $HOME/.jzero/config.yaml)
  -h, --help            help for jzero
  -t, --toggle          Help message for toggle

Use "jzero [command] --help" for more information about a command.
```

## jzero 命令行补全

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

