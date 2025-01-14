---
title: serverless 插件文档
icon: catppuccin:serverless
order: 3.5
---

## 新建 core 核心模块

```shell
jzero new core --feature serverless_core
cd core
jzero gen
go mod tidy
```

## 新建业务模块

```shell
cd core
jzero new b1 --feature serverless --output plugins/b1
cd plugins/b1
jzero gen
go mod tidy
```

## 构建 serverless

```shell
cd core
jzero serverless build

go run main.go server
```