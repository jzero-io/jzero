---
title: 生成服务端代码
icon: vscode-icons:folder-type-api-opened
order: 4
---

jzero 根据可描述语言生成代码:
* desc/api
* desc/proto
* desc/sql

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

## 下载依赖

```shell
go mod tidy
```

## 运行项目

```shell
go run main.go server
```

## 高级教程

### 基于 git 变动生成代码

```shell
jzero gen --git-change
```

### 指定 desc 生成代码

```shell
jzero gen --desc desc/api/xx.api
jzero gen --desc desc/proto/v1/xx.proto
jzero gen --desc desc/sql/xx.sql
```

### 生成代码忽略 desc

> 支持传入数组, 支持指定文件夹或者文件

```shell
jzero gen --desc-ignore desc/api/xx.api
jzero gen --desc-ignore desc/proto/v1/xx.proto
jzero gen --desc-ignore desc/sql/xx.sql
```