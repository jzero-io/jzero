---
home: false
icon: home
title: 首页
---

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/jzero.jpg" style="width: 33%;" alt=""/>
</div>

## 特性

* 企业级代码规范
* 支持自定义模板, 基于模板新建项目和生成代码, 默认支持多场景开发模板
* 一键创建项目, 一键生成服务端代码, 数据库代码, 客户端代码
* 基于 go-zero 框架, 扩展 go-zero 功能，并能完全同步 go-zero 新特性

## 快速开始

:::tip Windows 用户请在 powershell 下执行所有指令
:::

::: code-tabs#shell

@tab Docker

```bash
# 一键创建项目
docker run --rm -v ${PWD}/your_project:/app/your_project jaronnie/jzero:latest new your_project
cd your_project 
# 一键生成代码
docker run --rm -v ${PWD}:/app/your_project jaronnie/jzero:latest gen -w your_project
# 下载依赖
go mod tidy
# 启动项目
go run main.go server
```

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
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
```
:::