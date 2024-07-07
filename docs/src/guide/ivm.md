---
title: 接口版本控制
icon: hugeicons:api
order: 6
category: 开发
tag:
  - Guide
---

## 说明

接口版本控制功能, 是用来管理服务端的接口版本, 目前仅支持 grpc 场景. 默认会创建 v1 版本, 对应 desc/proto/v1 文件夹中的 proto

可以通过 ivm 命令自动初始化 v2 版本的接口, 并默认调用 v1 接口逻辑, 这意味着你仅需一条命令, 就可以自动生成 v2 的接口, 后续对 v2 接口继续更改即可

```shell
$ tree desc 
desc
└── proto
    ├── v1
    │   └── hello.proto
    └── v2
        └── hello_v2.proto
```

> 当前不支持 proto stream 的情况

## 初始化新版本

> 请依次进行变更, 如初始化 v3 版本时必须已有 v2 版本

```shell
# 初始化 v2 版本
jzero ivm init --version v2
# 或 jzero ivm init --version v2 --remove-suffix
```

## 新增 proto

可基于该命令自动生成一个带版本的 proto example, 可以快速生成一个 proto

```shell
# 在 desc/proto/v2 文件夹新增一个 machine.proto, service 默认与 name 同名, 可以指定 services, 也可以指定 service methods
jzero ivm add proto --name machine --version v2
```