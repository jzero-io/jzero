---
title: 生成服务端代码
icon: vscode-icons:folder-type-api-opened
order: 4
---

jzero 根据可描述语言生成代码:
* desc/api: api 可描述语言, 生成 http 服务端/客户端代码. [使用指南](../guide/api.md)
* desc/proto: proto 可描述语言, 生成 grpc 服务端/客户端代码. [使用指南](../guide/proto.md)
* desc/sql: sql 可描述语言, 生成数据库代码. [使用指南](../guide/model.md)
* model datasource: 通过远程数据源生成数据库代码. [使用指南](../guide/model.md)
* mongo: 通过指定 mongo type 生成 mongo 代码. [使用指南](../guide/mongo.md)

## 生成代码

::: code-tabs#shell

@tab jzero

```bash
cd your_project
jzero gen
```

@tab Docker

```bash
cd your_project
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
```
:::

## 基于 git 变动生成代码

```shell
jzero gen --git-change
```

## 指定 desc 生成代码

```shell
jzero gen --desc desc/api/xx.api
jzero gen --desc desc/proto/xx.proto
jzero gen --desc desc/sql/xx.sql
```

## 忽略指定 desc 生成代码

```shell
jzero gen --desc-ignore desc/api/xx.api
jzero gen --desc-ignore desc/proto/xx.proto
jzero gen --desc-ignore desc/sql/xx.sql
```

更多用法请参阅: [jzero 指南](../guide/jzero.md)