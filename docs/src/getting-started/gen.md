---
title: 生成服务端代码
icon: vscode-icons:folder-type-api-opened
order: 4
---

jzero 生成代码命令极其精简, 仅需 `jzero gen` 就能自动识别所有的可描述文件/配置, 完成代码的生成.

通过上一篇文档的 `jzero add` 命令添加可描述文件后, 执行 `jzero gen` 即可看到生成的文件了.

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

::: tip 基于 git status -su 获取新增/改动的可描述文件
:::

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