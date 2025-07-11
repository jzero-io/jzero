---
title: serverless 插件化特性
icon: arcticons:game-plugins
star: true
order: 5.4
category: 开发
tag:
  - Guide
---

jzero 支持插件化机制, 可以方便的进行插件的安装和卸载操作.

## 新增插件

```shell
# api 项目插件
jzero new your_plugin --frame api --serverless

# rpc 项目插件
jzero new your_plugin --frame rpc --serverless

# gateway 项目插件
jzero new your_plugin --frame gateway --serverless
```

## 编译项目

```shell
jzero serverless build

go build
```

## 卸载插件

```shell
# 卸载所有
jzero serverless delete

# 卸载指定插件
jzero serverless delete --plugin <plugin-name>

# 重新编译
go build
```