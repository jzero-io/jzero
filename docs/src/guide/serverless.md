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

## 新增需要支持插件化机制的项目

```shell
jzero new your_project --frame api --features serverless_core

cd your_project
jzero gen

go mod tidy
```

## 新增插件

```shell
jzero new your_plugin --frame api --features serverless --output ./plugins/your_plugin
cd ./plugins/your_plugin
jzero gen

go mod tidy
```

## 编译带有插件的项目

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