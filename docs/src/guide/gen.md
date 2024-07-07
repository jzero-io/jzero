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

### 高级用法

```shell
# 删除多余的生成的文件中后缀带 logic 和 handler 和 server. 以及在文件内容中多余的 Logic 和 Handler 和 Server
jzero gen --remove--suffix
```

## 下载依赖

```shell
go mod tidy
```

## 运行项目

```shell
go run main.go server
```
