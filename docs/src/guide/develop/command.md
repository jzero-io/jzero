---
title: 新增 command
icon: puzzle-piece
star: true
order: 0
category: 开发
tag:
  - Guide
---

jzero 基于 [cobra](https://github.com/spf13/cobra) 库实现命令行管理. 可基于 [cobra-cli](https://github.com/spf13/cobra-cli) 工具新增 command.

```shell
go install github.com/spf13/cobra-cli@latest

cd app1
cobra-cli add init

go run main.go -h

$ go run main.go -h
app1 root.

Usage:
  app1 [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  daemon      app1 daemon
  help        Help about any command
  init        A brief description of your command
  version     app1 version

Flags:
      --config string   config file (default is $HOME/.app1/config.yaml)
  -h, --help            help for app1
  -t, --toggle          Help message for toggle

Use "app1 [command] --help" for more information about a command.
```