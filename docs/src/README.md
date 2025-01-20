---
home: false
icon: fluent:home-heart-20-filled
title: 首页
---

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/jzero.jpg" style="width: 33%;" alt=""/>
</div>

## 介绍

基于 [go-zero](https://go-zero.dev) 开发的低代码微服务开发框架 `jzero`, 通过可描述文件(**api/proto/sql**)自动生成**服务端代码/客户端代码/数据库**代码, 降低开发心智, 解放双手!



jzero 具备以下特点:

* 支持通过配置文件, 命令行参数以及环境变量的组合的方式控制命令的参数, 告别繁琐的命令配置
* 支持基于 git 对改动文件部分生成代码, 极大提升大型项目代码生成效率
* 优化 go-zero 已有的痛点并扩展新的特性
* 内置常用开发模板并增强模板特性, 支持通过自定义模板内容, 构建企业内部代码模板
* 所有配套工具链跨平台使用, 支持 windows/mac/linux

## 快速开始

:::tip Windows 用户请在 powershell 下执行所有指令
:::

::: code-tabs#shell
@tab jzero

```bash
# 安装 jzero
go install github.com/jzero-io/jzero@latest
# 一键安装所需的工具
jzero check
# 一键创建项目
jzero new your_project
cd your_project
# 一键生成代码
jzero gen
# 生成 swagger
jzero gen swagger
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```

@tab Docker

```bash
# 创建项目
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest new your_project
cd your_project 
# 生成代码
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen
# 生成 swagger
docker run --rm -v ${PWD}:/app ghcr.io/jzero-io/jzero:latest gen swagger
# 下载依赖
go mod tidy
# 启动项目
go run main.go server
# 访问 swagger ui
http://localhost:8001/swagger
```
:::



