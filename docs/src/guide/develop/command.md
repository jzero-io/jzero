---
title: 命令行 cobra 使用文档
icon: clarity:command-line
star: true
order: 0.1
category: 开发
tag:
  - Guide
---

jzero 基于 [cobra](https://github.com/spf13/cobra) 库实现命令行管理. 可基于 [cobra-cli](https://github.com/spf13/cobra-cli) 工具新增 command.

```shell
go install github.com/spf13/cobra-cli@latest

cd your_project
cobra-cli add init
go run main.go -h
```