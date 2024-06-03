---
home: false
icon: home
title: 首页
---

可支持任意框架的脚手架 jzero, 默认支持 go-zero

::: tip 目前还在定制规范中, 不能保证兼容性. 待 v1.0 后保证兼容性
:::

<div style="text-align: center;">
  <img src="https://oss.jaronnie.com/jzero.jpg" style="width: 33%;" alt=""/>
</div>

## 特性

* 企业级代码规范
* grpc, grpc-gateway, api 三合一, 满足绝大部分场景业务需要
* 集成命令行框架 cobra, 轻松编写具备生产可用的命令行工具
* 支持多 proto 多 service, 减少开发耦合性
* 不修改源码, 完全同步 go-zero 新特性
* 一键创建项目, 快速拓展新业务, 减少心理负担
* 一键生成服务端代码, 数据库代码, 客户端 sdk, 大大提高开发测试效率
* 支持自定义模板, 基于模板新建项目和生成代码

## 快速开始

![2024-04-30_10-10-52](https://oss.jaronnie.com/2024-04-30_10-10-52.gif)

:::tip Windows 用户请在 powershell 下执行所有指令
:::

::: code-tabs#shell

@tab Docker

```bash
# 一键创建项目
docker run --rm -v ${PWD}/quickstart:/app/quickstart jaronnie/jzero:latest new quickstart
cd quickstart 
# 一键生成代码
docker run --rm -v ${PWD}:/app/quickstart jaronnie/jzero:latest gen -w quickstart
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
jzero new quickstart
cd quickstart
# 一键生成代码
jzero gen
# 下载依赖
go mod tidy
# 启动服务端程序
go run main.go server
```
:::