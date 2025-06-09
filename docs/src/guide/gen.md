---
title: 生成服务端代码
icon: vscode-icons:folder-type-api-opened
order: 4
---

jzero 根据可描述语言生成代码:
* desc/api: api 可描述语言, 生成 http 服务端/客户端代码. [使用文档](develop/api.md)
* desc/proto: proto 可描述语言, 生成 grpc 服务端/客户端代码. [使用文档](develop/proto.md)
* desc/sql: sql 可描述语言, 生成数据库代码. [使用文档](develop/model.md)
* model datasource: 通过远程数据源生成数据库代码. [使用文档](develop/model.md)

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

### 不同姿势使用 jzero

:::important 涨知识的小技巧
:::

* 支持通过配置文件 .jzero.yaml 控制各种参数
* 支持通过 flag 控制各种参数
* 支持通过环境变量控制各种参数
* 支持通过以上组合的方式控制各种参数, 优先级从高到低为 环境变量  > flag  > 配置文件

如: `jzero gen --style go_zero` 对应 .jzero.yaml 内容

```yaml
gen:
  style: go_zero
```

即 `jzero gen` + `.jzero.yaml` = `jzero gen --style go_zero`

对于环境变量的使用, 需要增加前缀 `JZERO_`, 如 `JZERO_GEN_STYLE`

即 `JZERO_GEN_STYLE=go_zero jzero gen` = `jzero gen --style go_zero`

### 基于 git 变动生成代码

```shell
jzero gen --git-change
```

```yaml
gen:
  git-change: true
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