---
home: false
icon: fluent:home-heart-20-filled
title: 首页
---

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/hiking.svg" style="width: 33%;" alt=""/>
</div>

## 介绍

一键新增 api/gateway/rpc 项目, 并基于可描述文件(**api/proto/sql**)自动生成**服务端和客户端代码**代码, 降低开发心智, 解放双手!

具备以下特点:

* 支持通过配置文件/命令行参数/环境变量组合的方式灵活控制 jzero 的各项配置, 极简指令生成代码, ai 友好
* 支持基于 git 对改动文件部分生成代码, 支持对指定描述文件生成代码或忽略指定描述文件生成代码, 提升大型项目代码生成效率
* 内置常用开发模板并增强模板特性, 支持自定义模板, 构建专属企业内部代码模板, 极大降低开发成本

更多详情请参阅：https://docs.jzero.io

## 快速开始

:::tip Windows 用户请在 powershell 下执行所有指令
:::

::: code-tabs#shell
@tab jzero

```bash
# 安装 jzero
go install github.com/jzero-io/jzero/cmd/jzero@latest
# 一键安装所需的工具
jzero check
# 一键创建项目
jzero new your_project
cd your_project
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```

@tab Docker

```bash
# 一键创建项目
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project 
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```
:::



