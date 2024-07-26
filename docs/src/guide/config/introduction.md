---
title: 配置说明
icon: vscode-icons:file-type-gleamconfig
order: 0
category: 配置
tag:
  - Guide
---

# 配置项

jzero 支持以下几种场景:

* api 项目: 基于 go-zero api 框架, 配置项为 **rest**
* rpc 项目: 基于 go-zero zrpc 框架, 配置项为 **zrpc**
* gateway 项目: 基于 go-zero gateway 框架以及 zrpc 框架, 配置项为 **gateway** 和 **zrpc**

其中对于基础服务配置项为 **rest**/**zrpc**/**gateway**, 日志配置项为 **log**

:::tip
后续所有关于配置的地方, 标题表示在哪个配置项进行配置
请根据自己的场景选取在哪个配置项配置
:::