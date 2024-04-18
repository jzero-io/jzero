---
title: proto 规范
icon: puzzle-piece
star: true
order: 1
category: 开发
tag:
  - Guide
---

:::tip jzero 支持多 proto 进行管理 proto(goctl 原生工具不支持).

jzero 在自动生成代码的时候会自动识别 daemon/proto 下的文件并自动注册到 zrpc 上.
:::

jzero 框架的理念是:

* 不同模块分在不同的 proto 文件下. 如一个系统, 用户模块即 user.proto, 权限模块即 auth.proto

proto 规范:

* 依据于 go-zero 的 proto 规范. 即 service 的 rpc 方法中 入参和出参的 proto 不能是 import 的 proto 文件