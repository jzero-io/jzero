---
title: 生成服务端代码
icon: vscode-icons:folder-type-api-opened
order: 4
---

jzero gen 根据 desc/api, desc/proto, desc/sql 文件夹下的文件生成代码. 生成代码的逻辑是调用 goctl 工具完成.

jzero 会自动检测对应文件夹下的内容, 然后进行自动生成, 使用者无需关心.

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

### 生成数据库代码

:::tip 支持基于 sql 文件和数据库 dsn 连接生成代码, 默认使用 sql 文件生成

[model 模板地址](https://github.com/jzero-io/sqlbuilder-zero) 

[go-zero model 文档](https://go-zero.dev/docs/tutorials/cli/model#goctl-model-mysql-%E6%8C%87%E4%BB%A4)
:::

jzero 默认使用 go-zero sqlx 和 sqlbuilder-go 完成对数据库的 crud 操作.

jzero 数据库规范:

* sql 文件放在 desc/sql
* 生成的 model 放在 internal/model

#### 基于 sql 文件生成代码

将 sql 文件放入 `desc/sql` 文件夹下即可, 执行 `jzero gen`

#### 基于数据库 dsn 连接生成代码

```shell
jzero gen --model-mysql-datasource --model-mysql-datasource-url="root:123456@tcp(127.0.0.1:3306)/test"
```

### 基于配置文件生成代码

基于[配置文件](./jzero.md#基于配置文件使用-jzero)生成服务端代码, 支持设置 before 和 after hooks

```yaml
syntax: v1

gen:
  hooks:
    before:
      - go run tools/migrate/main.go
    after:
      - jzero gen swagger

  # 从数据库连接生成 model 代码
  model-mysql-datasource: true
  model-mysql-datasource-url: "root:123456@tcp(127.0.0.1:3306)/test"
  model-mysql-ignore-columns: []
```
