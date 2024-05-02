---
title: 新增路由
icon: puzzle-piece
star: true
order: 1
category: 开发
tag:
  - Guide
---

jzero 提供了 3 种方式添加路由:

* grpc 与 grpc-gateway. 一套协议提供 grpc 与 http 接口(首选).
* go-zero api. 第一种方式不能满足业务的场景下使用, 如与文件相关的服务等.
* 自定义 route. 前两种无法满足业务场景下使用, 如将前端项目嵌入后端场景.