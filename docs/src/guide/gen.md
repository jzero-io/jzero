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
docker run --rm -v ${PWD}:/app jaronnie/jzero:latest gen
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

基于[配置文件](./command.md#基于配置文件使用-jzero)生成服务端代码, 支持设置 before 和 after hooks

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
